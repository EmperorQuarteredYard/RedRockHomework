package services

import (
	"fmt"
	controllers2 "lesson5/backend/controllers"
	"lesson5/backend/utils/JWT"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// PublicService 提供基础服务(登录,/static文件)
type PublicService struct {
	router         *gin.Engine
	loginHtml      []byte
	studentUser    *controllers2.StudentUser
	administerUser *controllers2.AdministerUser
}

// NewPublicService 创建并返回一个新的BaseService对象
func NewPublicService(studentUser *controllers2.StudentUser, administerUser *controllers2.AdministerUser, router *gin.Engine) *PublicService {
	var bs = PublicService{
		studentUser:    studentUser,
		administerUser: administerUser,
		router:         router,
	}
	bs.initHtmlFile()
	return &bs
}

// initHtmlFile 初始化登录的Html文件
func (s *PublicService) initHtmlFile() {
	s.loginHtml = loadHtml("login")
}

func (s *PublicService) Serve(wait *sync.WaitGroup) {
	fmt.Println("base server running on localhost:8080")
	s.RegisterRoutes()
	s.router.Run("localhost:8080")
	wait.Done()
}

func (s *PublicService) RegisterRoutes() {
	s.registerStaticFilesRoutes()
	s.registerPublicRoutes()

}

// registerStaticFilesRoutes /static路径初始化
func (s *PublicService) registerStaticFilesRoutes() {
	s.router.Static("/static", "C:/CodingProjects/RedRockHomework/Back-end/lesson5/frontend/static")
}

// registerPublicRoutes 注册登录与JWT路由
func (s *PublicService) registerPublicRoutes() {

	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})
	// 登录页面
	s.router.GET("/login", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", s.loginHtml)
	})

	// 登录API:
	s.router.POST("/api/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		role, err := s.administerUser.Login(username, password)
		if err != nil || role == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// 根据用户名获取用户ID（这里需要根据实际情况调整）
		userID := uint64(1) // 临时值，实际应从数据库获取

		// 生成JWT token
		roleStr := "student"
		if role > 1 {
			roleStr = "admin"
		}

		accessToken, refreshToken, err := lessonSelectJWT.GenerateToken(userID, roleStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"role":          roleStr,
			"redirect_url":  s.getRedirectURL(role),
		})
	})

	s.router.POST("/api/refresh", func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		claims, err := lessonSelectJWT.VerifyRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		newAccessToken, newRefreshToken, err := lessonSelectJWT.GenerateToken(claims.UserID, claims.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	})
}

// getRedirectURL 根据权限等级给出重定向的URL
func (s *PublicService) getRedirectURL(role int64) string {
	if role > 1 {
		return "/admin/lessons"
	}
	return "/student/lessons"
}
