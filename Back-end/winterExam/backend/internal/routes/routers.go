package routes

import (
	"homeworkSystem/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s *service.Service) {
	//auth := r.Group("/api/auth")
	//{
	//	auth.POST("/register", handlers.Register(db))
	//	auth.POST("/login", handlers.Login(db))
	//}
	//
	//// 需要认证的路由
	//api := r.Group("/api")
	//api.Use(middleware.AuthMiddleware())
	//{
	//	// 身份管理
	//	profiles := api.Group("/profiles")
	//	{
	//		profiles.GET("", handlers.GetProfiles(db))
	//		profiles.POST("", handlers.CreateProfile(db))
	//		profiles.PUT("/:id", handlers.UpdateProfile(db))
	//		profiles.DELETE("/:id", handlers.DeleteProfile(db))
	//	}
	//	// 密码条目管理
	//	entries := api.Group("/entries")
	//	{
	//		entries.GET("", handlers.GetEntries(db))
	//		entries.POST("", handlers.CreateEntry(db))
	//		entries.PUT("/:id", handlers.UpdateEntry(db))
	//		entries.DELETE("/:id", handlers.DeleteEntry(db))
	//	}
	//}
}
