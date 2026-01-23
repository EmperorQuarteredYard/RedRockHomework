package services

import (
	"encoding/json"
	"lesson4/controllers"
	"lesson4/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type StudentService struct {
	loginHtml               []byte
	lessonEditHtml          []byte
	selectionForStudentHtml []byte
	studentEditHtml         []byte
	studentLessonHtml       []byte
	selectionEditHtml       []byte
	initHtmlFileEvent       sync.Once
	router                  *gin.Engine
	studentUser             *controllers.StudentUser
}

func NewStudentService(studentUser *controllers.StudentUser) *StudentService {
	return &StudentService{studentUser: studentUser}
}

func (s *StudentService) initHtmlFile() {
	s.initHtmlFileEvent.Do(func() {
		s.loginHtml = loadHtml("login")
		s.lessonEditHtml = loadHtml("lessonEdit")
		s.selectionForStudentHtml = loadHtml("selectForStudent")
		s.studentEditHtml = loadHtml("studentEdit")
		s.selectionEditHtml = loadHtml("selectionEdit")
	})
}

func (s *StudentService) Serve() {
	s.initHtmlFile()
	s.registerRoutes()
	s.router.Run(":8080")
}

func (s *StudentService) registerRoutes() {
	s.router = gin.Default()
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})
	s.registerStudentSelectionRoutes()
	s.registerStaticFilesRoutes()
}

func (s *StudentService) registerStudentSelectionRoutes() {

	s.router.GET("/selectLesson", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.selectionForStudentHtml)
	})
	s.router.POST("/selectLesson", func(c *gin.Context) {
		lessons, err := s.studentUser.ListLessons()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"data": *lessons})
	})
	s.router.POST("/selectLesson/drop", func(c *gin.Context) {
		var req struct {
			StudentID int64 `json:"student_id"`
			LessonID  int64 `json:"lesson_id"`
		}
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
		}
		success, err := s.studentUser.Drop(req.StudentID, req.LessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can not drop"})
		}

	})
	s.router.POST("/selectLesson/select", func(c *gin.Context) {
		var req struct {
			StudentID int64 `json:"student_id"`
			LessonID  int64 `json:"lesson_id"`
		}
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
		}
		success, err := s.studentUser.Select(req.StudentID, req.LessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can not select"})
		}
	})
	s.router.POST("/selectLesson/check", func(c *gin.Context) {
		var req struct {
			StudentID int64 `json:"student_id"`
		}
		var lessons *[]models.Lesson
		lessons, err := s.studentUser.GetCurrentSelectedLesson(req.StudentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"data": *lessons})
	})
}

func (s *StudentService) registerStaticFilesRoutes() {
	s.router.Static("/static", "C:/CodingProjects/RedRockHomework/Back-end/lesson4/static")
}
