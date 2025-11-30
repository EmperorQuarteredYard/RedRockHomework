package services

import (
	"errors"
	"student_system/database"
	"student_system/models"
	"student_system/utils"
)

type AuthService struct{}

// 用户注册 - 修复返回值作用域问题
func (s *AuthService) Register(username, password, email string) error {
	var existingUser models.USER
	// 使用 = 而不是 := 来避免重新声明 err
	if err := database.DB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return errors.New("用户名或邮箱已存在")
	}

	user := models.USER{
		Username: username,
		Email:    email,
		Role:     "user",
	}

	if err := user.HashPassword(password); err != nil {
		return err
	}

	return database.DB.Create(&user).Error
}

// 用户登录 - 修复返回值作用域问题
func (s *AuthService) Login(username, password string) (string, *models.USER, error) {
	var user models.USER
	// 使用明确的错误处理
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", nil, errors.New("用户不存在")
	}

	if err := user.CheckPassword(password); err != nil {
		return "", nil, errors.New("密码错误")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", nil, errors.New("Token生成失败")
	}

	return token, &user, nil
}

// 获取当前用户信息
func (s *AuthService) GetCurrentUser(userID int) (*models.USER, error) {
	var user models.USER
	err := database.DB.First(&user, userID).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// 创建管理员用户
func (s *AuthService) CreateAdminUser() error {
	var adminUser models.USER
	// 检查管理员是否已存在
	err := database.DB.Where("username = ?", "admin").First(&adminUser).Error
	if err != nil {
		return err // 管理员已存在
	}

	adminUser = models.USER{
		Username: "admin",
		Email:    "admin@student.com",
		Role:     "admin",
	}

	if err := adminUser.HashPassword("admin123"); err != nil {
		return err
	}

	return database.DB.Create(&adminUser).Error
}
