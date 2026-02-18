package router

import (
	"homeworkSystem/backend/internal/controller"
	"homeworkSystem/backend/internal/service"
	"homeworkSystem/backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 初始化 service
	svc := service.NewService(db)

	// 初始化 controller
	userCtl := controller.NewUserController(svc)
	homeworkCtl := controller.NewHomeworkController(svc)
	submissionCtl := controller.NewSubmissionController(svc)

	// 公开接口
	public := r.Group("/api")
	{
		public.POST("/user/register", userCtl.Register)
		public.POST("/user/login", userCtl.Login)
		public.POST("/user/refresh", userCtl.RefreshToken)
	}

	// 需要认证的接口
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuth())
	{
		// 用户模块
		auth.GET("/user/profile", userCtl.GetProfile)
		auth.DELETE("/user/account", userCtl.DeleteAccount)

		// 作业模块
		auth.POST("/homework", homeworkCtl.Publish)
		auth.GET("/homework", homeworkCtl.List)
		auth.GET("/homework/:id", homeworkCtl.GetDetail)
		auth.PUT("/homework/:id", homeworkCtl.Update)
		auth.DELETE("/homework/:id", homeworkCtl.Delete)

		// 提交模块
		auth.POST("/submission", submissionCtl.Submit)
		auth.GET("/submission/my", submissionCtl.MySubmissions)
		auth.GET("/submission/homework/:homework_id", submissionCtl.HomeworkSubmissions)
		auth.PUT("/submission/:id/review", submissionCtl.Review)
		auth.PUT("/submission/:id/excellent", submissionCtl.MarkExcellent)
		auth.GET("/submission/excellent", submissionCtl.ExcellentList)
	}

	return r
}
