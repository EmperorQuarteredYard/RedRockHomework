package services

import (
	"encoding/json"
	"fmt"
	"lesson4/controllers"
	"lesson4/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdministerService struct {
	*StudentService
	AdministerUser *controllers.AdministerUser
}

func NewAdministerService(administerUser *controllers.AdministerUser, studentUser *controllers.StudentUser) *AdministerService {
	return &AdministerService{AdministerUser: administerUser, StudentService: NewStudentService(studentUser)}
}

func (s *AdministerService) Serve() {
	s.initHtmlFile()
	s.StudentService.registerRoutes()
	s.registerRoutes()
	fmt.Println("serve run on localhost:8080")
	s.router.Run("localhost:8080")
}

func (a *AdministerService) registerRoutes() {
	a.StudentService.initHtmlFile()
	a.StudentService.registerRoutes()
	a.beginLoginServe()
	a.beginSelectionEditServe()
	a.beginStudentEditServe()
	a.beginLessonEditServe()
}

func (s *AdministerService) beginLoginServe() {
	s.router.GET("/login", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.loginHtml)
	})
	s.router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		fmt.Println("用户登入:")
		fmt.Println("账户:" + username)
		fmt.Println("密码:" + password)
		role, err := s.AdministerUser.Login(username, password)
		fmt.Println("权限:", role)
		if err != nil || role == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无权访问",
			})
			return
		}

		if role > 1 {
			c.Redirect(http.StatusMovedPermanently, "/lessonEdit")
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/selectLesson")
	})
}

func (s *AdministerService) beginLessonEditServe() {
	s.router.GET("/lessonEdit", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.lessonEditHtml)
	})
	s.router.POST("/lessonEdit", func(c *gin.Context) {
		lessons, _ := s.AdministerUser.ListLessons()
		c.JSON(200, gin.H{
			"lessons": lessons,
		})
	})
	s.router.POST("/lessonEdit/create", func(c *gin.Context) {
		var req models.Lesson

		// 读取原始请求体
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
		}
		err = s.AdministerUser.CreateLesson(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{})
	})
	s.router.POST("/lessonEdit/delete", func(c *gin.Context) {
		var req struct {
			ID int64
		}
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			fmt.Println("无法读取请求体")
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
			fmt.Println("无法解析请求体")
			return
		}
		err = s.AdministerUser.DeleteLesson(req.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println(err.Error())
			return
		}
	})
}

func (s *AdministerService) beginStudentEditServe() {
	s.router.GET("/studentEdit", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.studentEditHtml)
	})
	s.router.POST("/studentEdit", func(c *gin.Context) {
		students, _ := s.AdministerUser.ListStudent()
		fmt.Println(students)
		c.JSON(200, gin.H{
			"students": students,
		})
	})
	s.router.POST("/studentEdit/create", func(c *gin.Context) {
		var req models.Student

		// 读取原始请求体
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			fmt.Println("无法读取请求体")
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
			fmt.Println("无法解析请求体")
			return
		}
		fmt.Println("获取学生", req)
		err = s.AdministerUser.CreateStudent(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println(err.Error())
			return
		}
		c.JSON(200, gin.H{})
		fmt.Println("存储学生", req)
	})
	s.router.POST("/studentEdit/delete", func(c *gin.Context) {
		var req struct {
			ID int64
		}
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
			fmt.Println("无法解析请求体")
			return
		}
		fmt.Println("删除学生", req.ID)
		err = s.AdministerUser.DeleteStudent(req.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println(err.Error())
			return
		}
	})
}

func (s *AdministerService) beginSelectionEditServe() {
	s.router.GET("/selectionEdit", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.selectionEditHtml)
	})
	s.router.POST("/selectionEdit", func(c *gin.Context) {
		selections, _ := s.AdministerUser.ListSelection()
		c.JSON(200, gin.H{
			"selections": selections,
		})
	})
	s.router.POST("/selectionEdit/create", func(c *gin.Context) {
		var req models.Selection

		// 读取原始请求体
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析请求体"})
		}
		success, err := s.AdministerUser.Select(req.StudentID, req.LessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can not select"})
		}
		c.JSON(200, gin.H{})
	})
	s.router.POST("/selectionEdit/delete", func(c *gin.Context) {
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
		success, err := s.AdministerUser.Drop(req.StudentID, req.LessonID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can not drop"})
		}
	})
}
