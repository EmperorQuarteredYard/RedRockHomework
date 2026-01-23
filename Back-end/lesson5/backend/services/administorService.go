// adminService.go
package services

import (
	"fmt"
	controllers2 "lesson5/backend/controllers"
	"lesson5/backend/models"
	"lesson5/backend/utils/JWT"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type AdminService struct {
	router            *gin.Engine
	lessonEditHtml    []byte
	studentEditHtml   []byte
	selectionEditHtml []byte
	studentUser       *controllers2.StudentUser
	administerUser    *controllers2.AdministerUser
}

func NewAdminService(studentUser *controllers2.StudentUser, administerUser *controllers2.AdministerUser, router *gin.Engine) *AdminService {

	service := &AdminService{
		router:         router,
		studentUser:    studentUser,
		administerUser: administerUser,
	}
	service.initAdminHtml()
	return service
}

func (s *AdminService) RegisterRoutes() {
	s.registerAdminRoutes()
}

func (s *AdminService) Serve(wait *sync.WaitGroup) {
	s.RegisterRoutes()
	fmt.Println("AdminService server running on localhost:8080")
	s.router.Run("localhost:8080")
	wait.Done()
}

func (s *AdminService) initAdminHtml() {
	s.lessonEditHtml = loadHtml("lessonEdit")
	s.studentEditHtml = loadHtml("studentEdit")
	s.selectionEditHtml = loadHtml("selectionEdit")
}

func (s *AdminService) registerAdminRoutes() {
	// 管理员路由组
	adminGroup := s.router.Group("/admin")
	adminGroup.Use(lessonSelectJWT.JWTAuthMiddleware())
	adminGroup.Use(lessonSelectJWT.RoleMiddleware("admin"))

	// 课程管理
	s.registerLessonRoutes(adminGroup)

	// 学生管理
	s.registerStudentRoutes(adminGroup)

	// 选课管理
	s.registerSelectionRoutes(adminGroup)
}

func (s *AdminService) registerLessonRoutes(adminGroup *gin.RouterGroup) {
	// 课程管理页面
	adminGroup.GET("/lessons", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.lessonEditHtml)
	})

	// 获取所有课程
	adminGroup.GET("/api/lessons", func(c *gin.Context) {
		lessons, err := s.administerUser.ListLessons()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"lessons": lessons})
	})

	// 创建课程
	adminGroup.POST("/api/lessons", func(c *gin.Context) {
		var req models.Lesson

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := s.administerUser.CreateLesson(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	})

	// 删除课程
	adminGroup.DELETE("/api/lessons/:id", func(c *gin.Context) {

		id := c.Param("id")
		var idInt int64
		fmt.Sscanf(id, "%d", &idInt)

		err := s.administerUser.DeleteLesson(idInt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	})
}

func (s *AdminService) registerStudentRoutes(adminGroup *gin.RouterGroup) {
	// 学生管理页面
	adminGroup.GET("/students", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.studentEditHtml)
	})

	// 获取所有学生
	adminGroup.GET("/api/students", func(c *gin.Context) {
		students, err := s.administerUser.ListStudent()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"students": students})
	})

	// 创建学生
	adminGroup.POST("/api/students", func(c *gin.Context) {
		var req models.Student

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := s.administerUser.CreateStudent(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	})

	// 删除学生
	adminGroup.DELETE("/api/students/:id", func(c *gin.Context) {
		id := c.Param("id")
		var idInt int64
		fmt.Sscanf(id, "%d", &idInt)

		err := s.administerUser.DeleteStudent(idInt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": true})
	})
}

func (s *AdminService) registerSelectionRoutes(adminGroup *gin.RouterGroup) {
	// 选课管理页面
	adminGroup.GET("/selections", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.selectionEditHtml)
	})

	// 获取所有选课记录
	adminGroup.GET("/api/selections", func(c *gin.Context) {
		selections, err := s.administerUser.ListSelection()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"selections": selections})
	})

	// 创建选课记录
	adminGroup.POST("/api/selections", func(c *gin.Context) {
		var req models.Selection

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		success, err := s.administerUser.Select(req.StudentID, req.LessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create selection"})
			return
		}

		c.JSON(200, gin.H{"success": true})
	})

	// 删除选课记录
	adminGroup.DELETE("/api/selections", func(c *gin.Context) {
		var req struct {
			StudentID int64 `json:"student_id"`
			LessonID  int64 `json:"lesson_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		success, err := s.administerUser.Drop(req.StudentID, req.LessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete selection"})
			return
		}

		c.JSON(200, gin.H{"success": true})
	})
}
