// studentService.go
package services

import (
	"fmt"
	controllers2 "lesson5/backend/controllers"
	"lesson5/backend/utils/JWT"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type StudentService struct {
	router                  *gin.Engine
	selectionForStudentHtml []byte
	studentLessonHtml       []byte
	studentUser             *controllers2.StudentUser
	administerUser          *controllers2.AdministerUser
}

func NewStudentService(studentUser *controllers2.StudentUser, administerUser *controllers2.AdministerUser, router *gin.Engine) *StudentService {
	service := &StudentService{
		router:         router,
		studentUser:    studentUser,
		administerUser: administerUser,
	}
	service.initStudentHTMLFiles()
	return service
}

func (s *StudentService) initStudentHTMLFiles() {
	s.selectionForStudentHtml = loadHtml("selectForStudent")
	s.studentLessonHtml = loadHtml("studentLesson")
}

func (s *StudentService) Serve(wait *sync.WaitGroup) {
	s.RegisterRoutes() // 注册所有路由
	fmt.Println("StudentService server running on localhost:8080")
	s.router.Run("localhost:8080")
	wait.Done()
}
func (s *StudentService) RegisterRoutes() {
	// 注册学生特有的路由
	s.registerStudentRoutes()
}

func (s *StudentService) registerStudentRoutes() {
	// 学生路由组
	studentGroup := s.router.Group("/student")
	studentGroup.Use(lessonSelectJWT.JWTAuthMiddleware())
	studentGroup.Use(lessonSelectJWT.RoleMiddleware("student"))

	// 学生选课页面
	studentGroup.GET("/lessons", func(c *gin.Context) {
		s.initStudentHTMLFiles()
		c.Data(200, "text/html; charset=utf-8", s.selectionForStudentHtml)
	})

	// 我的课程页面
	studentGroup.GET("/my-lessons", func(c *gin.Context) {
		s.initStudentHTMLFiles()
		c.Data(200, "text/html; charset=utf-8", s.studentLessonHtml)
	})

	// API 路由
	apiGroup := studentGroup.Group("/api")
	{
		// 获取所有课程
		apiGroup.GET("/lessons", func(c *gin.Context) {
			lessons, err := s.studentUser.ListLessons()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": *lessons})
		})

		// 获取已选课程
		apiGroup.GET("/lessons/selected", func(c *gin.Context) {
			user, err := lessonSelectJWT.GetUserFromContext(c)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			lessons, err := s.studentUser.GetCurrentSelectedLesson(int64(user.UserID))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": *lessons})
		})

		// 选择课程
		apiGroup.POST("/lessons/select", func(c *gin.Context) {
			var req struct {
				StudentID int64 `json:"student_id"`
				LessonID  int64 `json:"lesson_id"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			success, err := s.studentUser.Select(req.StudentID, req.LessonID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if !success {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot select the lesson"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		// 退选课程
		apiGroup.POST("/lessons/drop", func(c *gin.Context) {
			var req struct {
				StudentID int64 `json:"student_id"`
				LessonID  int64 `json:"lesson_id"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			success, err := s.studentUser.Drop(req.StudentID, req.LessonID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if !success {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot drop the lesson"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true})
		})
	}
}
