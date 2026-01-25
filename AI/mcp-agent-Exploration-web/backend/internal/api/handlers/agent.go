package handlers

import (
	"backend/internal/agent"
	"backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AgentHandler Agent 处理器
type AgentHandler struct {
	agent  *agent.LLMAgent
	logger *zap.Logger
}

// NewAgentHandler 创建新的 Agent 处理器
func NewAgentHandler(agent *agent.LLMAgent, logger *zap.Logger) *AgentHandler {
	return &AgentHandler{
		agent:  agent,
		logger: logger,
	}
}

// GetStatus 获取 Agent 状态
// @Summary 获取Agent状态
// @Description 获取LLM Agent的当前状态和健康信息
// @Tags agent
// @Accept json
// @Produce json
// @Success 200 {object} models.AgentStatus
// @Router /api/status [get]
func (h *AgentHandler) GetStatus(c *gin.Context) {
	status := h.agent.GetStatus(c.Request.Context())

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    status,
	})
}

// GetTools 获取工具列表
// @Summary 获取可用工具列表
// @Description 获取所有可用的工具及其信息
// @Tags agent
// @Accept json
// @Produce json
// @Success 200 {array} models.ToolInfo
// @Router /api/tools [get]
func (h *AgentHandler) GetTools(c *gin.Context) {
	tools := h.agent.GetTools()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tools,
	})
}

// DecomposeTask 拆解任务
// @Summary 拆解任务
// @Description 将复杂任务拆解为子任务并估计时间
// @Tags task
// @Accept json
// @Produce json
// @Param request body models.TaskDecompositionRequest true "任务拆解请求"
// @Success 200 {object} models.TaskDecompositionResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/decompose [post]
func (h *AgentHandler) DecomposeTask(c *gin.Context) {
	var request models.TaskDecompositionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("请求参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 验证任务描述
	if len(request.TaskDescription) < 5 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "任务描述太短，请提供更多细节",
		})
		return
	}

	if len(request.TaskDescription) > 1000 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "任务描述过长，请简要描述",
		})
		return
	}

	h.logger.Info("处理任务拆解请求",
		zap.String("task", request.TaskDescription),
		zap.Int("max_steps", request.MaxSteps),
	)

	response, err := h.agent.DecomposeTask(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("任务拆解失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "任务拆解失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Description 根据任务ID获取任务拆解详情
// @Tags task
// @Accept json
// @Produce json
// @Param task_id path string true "任务ID"
// @Success 200 {object} models.TaskDecompositionResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/task/{task_id} [get]
func (h *AgentHandler) GetTask(c *gin.Context) {
	taskID := c.Param("task_id")

	task, exists := h.agent.GetTask(taskID)
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "任务不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    task,
	})
}

// QuickEstimate 快速时间估计
// @Summary 快速时间估计
// @Description 快速估计任务所需时间
// @Tags task
// @Accept json
// @Produce json
// @Param task query string true "任务描述"
// @Param complexity query string false "复杂度级别" Enums(auto, simple, medium, complex) default(medium)
// @Success 200 {object} map[string]interface{}
// @Router /api/estimate [post]
func (h *AgentHandler) QuickEstimate(c *gin.Context) {
	task := c.Query("task")
	complexity := c.DefaultQuery("complexity", "medium")

	if task == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "任务描述不能为空",
		})
		return
	}

	result := h.agent.QuickEstimate(c.Request.Context(), task, complexity)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}
