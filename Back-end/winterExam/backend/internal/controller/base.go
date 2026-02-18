package controller

import (
	"homeworkSystem/backend/internal/models"
	"homeworkSystem/backend/pkg/errcode"
	"homeworkSystem/backend/pkg/jwt"
	"homeworkSystem/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// 绑定 JSON 并自动处理错误
func (b *BaseController) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Error(c, errcode.ErrInvalidParams)
		return false
	}
	return true
}

// 获取当前认证用户
func (b *BaseController) GetAuthUser(c *gin.Context) (*jwt.AuthUser, bool) {
	val, exists := c.Get("user")
	if !exists {
		response.Error(c, errcode.ErrUnauthorized)
		return nil, false
	}
	user, ok := val.(jwt.AuthUser)
	if !ok {
		response.Error(c, errcode.ErrServer)
		return nil, false
	}
	return &user, true
}

// 获取部门标签
func (b *BaseController) DepartmentLabel(dept string) string {
	d, _ := models.GetDepartment(dept)
	return d.DepartmentLabel
}

// 成功响应（带数据）
func (b *BaseController) Success(c *gin.Context, data interface{}) {
	response.Success(c, data)
}

// 处理错误响应
func (b *BaseController) HandleError(c *gin.Context, err error) {
	response.Error(c, err)
}
