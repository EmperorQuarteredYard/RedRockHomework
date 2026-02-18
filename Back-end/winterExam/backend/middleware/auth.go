package middleware

import (
	"homeworkSystem/backend/pkg/jwt"
	"homeworkSystem/backend/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, jwt.ErrMissingToken)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, jwt.ErrInvalidToken)
			c.Abort()
			return
		}

		claims, err := jwt.VerifyAccessToken(parts[1])
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user", jwt.AuthUser{
			UserID:     claims.UserID,
			Department: claims.Department,
			Role:       claims.Role,
		})
		c.Next()
	}
}
