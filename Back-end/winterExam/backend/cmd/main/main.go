package main

import (
	"homeworkSystem/backend/configs"
	"homeworkSystem/backend/internal/router"
	"log"
)

func main() {
	// 加载配置
	cfg := configs.LoadConfig()

	// 初始化数据库
	db := configs.InitDB(cfg)

	// 初始化路由
	r := router.SetupRouter(db)

	// 启动服务
	if err := r.Run(":51443"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
