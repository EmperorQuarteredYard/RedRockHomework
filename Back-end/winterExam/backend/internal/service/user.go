package service

import (
	"homeworkSystem/backend/internal/models"
	WEjwt "homeworkSystem/backend/pkg/middleware/jwt"

	"github.com/gin-gonic/gin"
)

// Register 返回一个用于处理注册事件的gin.HandlerFunc，不需要jwtToken
func (s *Service) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var label string
		var message struct {
			Username   string `json:"username"`
			Password   string `json:"password"`
			Nickname   string `json:"nickname"`
			Department string `json:"department"`
		}
		delegation.handleDataBind(&message)
		delegation.handleDataForm(message.Username == "" || message.Password == "" || message.Nickname == "" || message.Department == "")
		delegation.handleLabelParse(message.Department, &label)

		if !delegation.Success() {
			return
		}
		user, err := s.userDao.Register(message.Username, message.Password, message.Department, message.Nickname, models.RoleNewLight)
		if err != nil {
			delegation.handleErrorAbort(err)
			return
		}
		delegation.handleSuccessResponse(gin.H{
			"id":               user.ID,
			"username":         user.Name,
			"nickname":         user.Nickname,
			"role":             user.Role,
			"department":       user.Department,
			"department_label": label,
		})
		return
	}
}

// Login 用户登录，不需要jwtToken
func (s *Service) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var user *models.User
		var err error
		var accessToken, refreshToken, label string
		var message struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		delegation.handleDataBind(&message)
		delegation.handleDataForm(message.Username == "" || message.Password == "")
		if !delegation.Success() {
			return
		}

		if user, err = s.userDao.Login(message.Username, message.Password); err != nil {
			delegation.handleErrorAbort(err)
			return
		}
		if accessToken, refreshToken, err = WEjwt.GenerateToken(user.ID, user.Department, user.Role); err != nil {
			delegation.handleErrorAbort(err)
			return
		}
		delegation.handleLabelParse(user.Department, &label)

		if delegation.Success() {
			delegation.handleSuccessResponse(
				gin.H{
					"access_token":  accessToken,
					"refresh_token": refreshToken,
					"user": gin.H{
						"id":               user.ID,
						"username":         user.Name,
						"nickname":         user.Nickname,
						"role":             user.Role,
						"department":       user.Department,
						"department_label": label,
					},
				})
		}
		return
	}
}

// RefreshToken 刷新 Token，head中不需要jwtToken
func (s *Service) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var err error
		var claim *WEjwt.CustomClaims
		var accessToken, refreshToken string

		var message struct {
			RefreshToken string `json:"refresh_token"`
		}
		delegation.handleDataBind(&message)
		delegation.handleDataForm(message.RefreshToken == "")
		if !delegation.Success() {
			return
		}

		if claim, err = WEjwt.VerifyRefreshToken(message.RefreshToken); err != nil {
			delegation.handleErrorAbort(err)
			return
		}

		if accessToken, refreshToken, err = WEjwt.GenerateToken(claim.UserID, claim.Department, claim.Role); err != nil {
			delegation.handleErrorAbort(err)
			return
		}

		delegation.handleSuccessResponse(gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
		return
	}
}

func (s *Service) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var authUser WEjwt.AuthUser
		var user *models.User
		var err error
		var label string

		delegation.handleUserVerify(&authUser)
		if user, err = s.userDao.FindByID(authUser.UserID); err != nil {
			delegation.handleErrorAbort(err)
			return
		}
		delegation.handleLabelParse(user.Department, &label)

		delegation.handleSuccessResponse(gin.H{
			"id":               user.ID,
			"username":         user.Name,
			"nickname":         user.Nickname,
			"role":             user.Role,
			"department":       user.Department,
			"department_label": label,
			"email":            user.Email,
		})
	}
}

func (s *Service) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		delegation := NewDelegation(c)
		var authUser WEjwt.AuthUser
		var message struct {
			Password string `json:"password"`
		}

		delegation.handleUserVerify(&authUser)
		delegation.handleDataBind(&message)
		delegation.handleDataForm(message.Password == "")
		if !delegation.Success() {
			return
		}

		err := s.userDao.SoftDelete(authUser.UserID, message.Password)
		if err != nil {
			delegation.handleErrorAbort(err)
		}
		delegation.handleSuccessResponse(nil)
	}
}
