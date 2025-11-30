package utils

import (
	"student_system/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // 生产环境请使用环境变量

// 生成JWT Token
func GenerateToken(user models.USER) (string, error) {
	now := time.Now()
	expireTime := now.Add(24 * time.Hour) // Token有效期24小时

	claims := models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "student_system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 验证JWT Token
func ParseToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// 从Token中获取用户ID
func GetUserIDFromToken(tokenString string) (int, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
