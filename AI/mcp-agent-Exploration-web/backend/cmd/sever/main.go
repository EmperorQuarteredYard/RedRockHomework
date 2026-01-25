package main

import (
	"backend/internal/agent"
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"
	"backend/internal/config"
	"backend/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var (
	version = "1.0.0"
	build   = "dev"
)

// @title LLM Agent API
// @version 1.0.0
// @description 基于本地LLM的任务规划助手API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("C:\\CodingProjects\\RedRockHomework\\AI\\mcp-agent-Exploration-web\\backend\\configs\\configs.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log, err := logger.NewLogger(cfg.Log.Level, cfg.Log.File)
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync(log)

	log.Info("启动 LLM Agent 后端服务",
		zap.String("version", version),
		zap.String("build", build),
		zap.String("mode", cfg.Server.Mode),
	)

	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 实例
	router := gin.New()

	// 注册中间件
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())

	// 创建 Agent 实例
	llmAgent := agent.NewLLMAgent(
		cfg.Ollama.ModelName,
		cfg.Ollama.BaseURL,
		log,
	)

	// 创建处理器
	healthHandler := handlers.NewHealthHandler(version)
	agentHandler := handlers.NewAgentHandler(llmAgent, log)
	chatHandler := handlers.NewChatHandler(llmAgent, log)

	// 注册路由
	setupRoutes(router, healthHandler, agentHandler, chatHandler)

	// 配置服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:  cfg.Server.IdleTimeout * time.Second,
	}

	// 启动服务器
	go func() {
		log.Info("服务器启动",
			zap.String("address", srv.Addr),
			zap.String("mode", cfg.Server.Mode),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("服务器关闭失败", zap.Error(err))
	}

	log.Info("服务器已关闭")
}

// setupRoutes 设置路由
func setupRoutes(
	router *gin.Engine,
	healthHandler *handlers.HealthHandler,
	agentHandler *handlers.AgentHandler,
	chatHandler *handlers.ChatHandler,
) {
	// 健康检查
	router.GET("/health", healthHandler.HealthCheck)

	// API 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := router.Group("/api")
	{
		// Agent 相关
		api.GET("/status", agentHandler.GetStatus)
		api.GET("/tools", agentHandler.GetTools)

		// 任务相关
		api.POST("/decompose", agentHandler.DecomposeTask)
		api.GET("/task/:task_id", agentHandler.GetTask)
		api.POST("/estimate", agentHandler.QuickEstimate)

		// 聊天相关
		api.POST("/chat", chatHandler.Chat)
	}

	// 根路径
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "LLM Agent API",
			"version": version,
			"docs":    "/swagger/index.html",
			"status":  "/health",
		})
	})
}
