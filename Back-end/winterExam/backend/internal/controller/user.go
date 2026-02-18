package controller

import (
	"github.com/gin-gonic/gin"
	"homeworkSystem/backend/internal/models"
	"homeworkSystem/backend/internal/service"
	"homeworkSystem/backend/pkg/errcode"
)

type UserController struct {
	BaseController
	svc *service.Service
}

func NewUserController(svc *service.Service) *UserController {
	return &UserController{svc: svc}
}

// Register 用户注册
func (ctl *UserController) Register(c *gin.Context) {
	var req struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		Nickname   string `json:"nickname" binding:"required"`
		Department string `json:"department" binding:"required"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	user, err := ctl.svc.Register(req.Username, req.Password, req.Nickname, req.Department)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"id":               user.ID,
		"username":         user.Username,
		"nickname":         user.Nickname,
		"role":             user.Role,
		"department":       user.Department,
		"department_label": ctl.DepartmentLabel(user.Department),
	})
}

// Login 登录
func (ctl *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	user, accessToken, refreshToken, err := ctl.svc.Login(req.Username, req.Password)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":               user.ID,
			"username":         user.Username,
			"nickname":         user.Nickname,
			"role":             user.Role,
			"department":       user.Department,
			"department_label": ctl.DepartmentLabel(user.Department),
		},
	})
}

// RefreshToken 刷新 token
func (ctl *UserController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	accessToken, refreshToken, err := ctl.svc.RefreshToken(req.RefreshToken)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetProfile 获取当前用户信息
func (ctl *UserController) GetProfile(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	user, err := ctl.svc.GetUserByID(authUser.UserID)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"id":               user.ID,
		"username":         user.Username,
		"nickname":         user.Nickname,
		"role":             user.Role,
		"department":       user.Department,
		"department_label": ctl.DepartmentLabel(user.Department),
		"email":            user.Email,
	})
}

// DeleteAccount 注销账号
func (ctl *UserController) DeleteAccount(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if !ctl.BindJSON(c, &req) {
		return
	}

	err := ctl.svc.DeleteAccount(authUser.UserID, req.Password)
	if err != nil {
		ctl.HandleError(c, err)
		return
	}

	ctl.Success(c, nil)
}
