package service

import (
	"fmt"
	"homeworkSystem/backend/internal/models"
	"homeworkSystem/backend/internal/repository"
	"homeworkSystem/backend/pkg/jwt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type GeneralDelegation struct {
	c       *gin.Context
	success bool
	lock    sync.Locker
}

// NewDelegation 构造一个新的委托
func NewDelegation(c *gin.Context) *GeneralDelegation {
	return &GeneralDelegation{c: c, success: true}
}

func (d *GeneralDelegation) handlePermissionRole(userRole, permittedRole string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if userRole == permittedRole {
		return
	}
	d.success = false
	d.c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"code":    StatusPermissionInsufficient,
		"message": repository.ErrorInsufficientPermissions,
		"data":    nil,
	})
}

// Success 之前的操作是否都成功了
func (d *GeneralDelegation) Success() bool {
	return d.success
}

// handleSuccessResponse 成功相应(摆出高雅人士的样子)
func (d *GeneralDelegation) handleSuccessResponse(data interface{}) {
	d.lock.Lock()
	if !d.success {
		d.lock.Unlock()
		return
	}
	d.c.JSON(http.StatusOK, gin.H{
		"code":    StatusSuccess,
		"message": "success",
		"data":    data,
	})
	d.lock.Unlock()
	return
}

// handleErrorAbort 大概报个错
func (d *GeneralDelegation) handleErrorAbort(err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if !d.success {
		return
	}
	d.success = false
	d.c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"code":    StatusUnknownError,
		"message": "something unknow went wrong",
		"data":    nil,
	})
	fmt.Println("something went wrong on database:" + err.Error())
	return
}

// handleDataForm 检查数据是否合法
func (d *GeneralDelegation) handleDataForm(condition bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if !d.success {
		return
	}
	if condition {
		d.success = false
		d.c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    StatusInvalidData,
			"message": "data invalid",
			"data":    nil,
		})
	}
	return
}

// handleDataBind 绑定数据(从body中)
func (d *GeneralDelegation) handleDataBind(message interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if !d.success {
		return
	}
	err := d.c.ShouldBindJSON(message)
	if err != nil {
		d.success = false
		d.c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    StatusDataMissed,
			"message": "fail to bind data",
			"data":    nil,
		})
		fmt.Println("fail to bind data" + err.Error())
	}
	return
}

// handleUserVerify 取得用户信息
func (d *GeneralDelegation) handleUserVerify(user *WEjwt.AuthUser) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if !d.success {
		return
	}
	bearerUser, ok := d.c.Get("user")
	if !ok {
		d.handleErrorAbort(fmt.Errorf("maybe the stupid programmer forget to use jwt middleware"))
		return
	}
	authUser, ok := bearerUser.(WEjwt.AuthUser)
	if !ok {
		d.handleErrorAbort(fmt.Errorf("user is not a user"))
		return
	}
	user = &authUser
	return
}

// handleLabelParse 获取部门的label
func (d *GeneralDelegation) handleLabelParse(departmentValue string, label *string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if !d.success {
		return
	}
	department, err := models.GetDepartment(departmentValue)
	if err != nil {
		d.success = false
		d.c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    StatusInvalidData,
			"message": "department invalid: " + departmentValue,
			"data":    nil,
		})
		return
	}
	s := fmt.Sprintf(department.DepartmentLabel)
	label = &s
	return
}
