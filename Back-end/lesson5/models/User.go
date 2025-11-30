package models

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type USER struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Role      string `gorm:"type:enum('admin','user');default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 哈希密码
func (u *USER) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// 验证密码
func (u *USER) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// JWT声明结构
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
