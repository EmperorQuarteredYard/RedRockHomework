// main.go
package main

import (
	"fmt"
	"lesson5/backend/controllers"
	"lesson5/backend/database"
	"lesson5/backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化控制器
	db := database.InitDB()
	if db == nil {
		fmt.Println("db为空")
	}
	studentUser := controllers.NewStudentUser(db)
	administerUser := controllers.NewAdministerUser(db)

	// 创建组合服务
	router := gin.Default()
	adminService := services.NewAdminService(studentUser, administerUser, router)
	studentService := services.NewStudentService(studentUser, administerUser, router)
	publicService := services.NewPublicService(studentUser, administerUser, router)

	go publicService.RegisterRoutes()
	go studentService.RegisterRoutes()
	go adminService.RegisterRoutes()

	router.Run(":8080")

	// 启动服务
}
