package DAO

import (
	"errors"
	"fmt"
	"lesson5/backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDAO struct {
	*BaseDAO
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	if db == nil {
		fmt.Println("NewUserDao Called: db不能为空")
		panic(errors.New("NewUserDao Called: db不能为空"))
		return nil
	}
	return &UserDAO{NewBaseDAO(db)}
}

func hash(s string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func (d *UserDAO) CreateUser(name, password, role string) error {
	var user = models.User{Name: name, PasswordHash: hash(password), Role: role}
	return d.db.Create(&user).Error
}

func (d *UserDAO) DeleteUserByID(userID int64) error {
	return d.db.Where("id = ?", userID).Delete(&models.User{}).Error
}
func (d *UserDAO) DeleteUserByName(name string) error {
	return d.db.Where("name = ?", name).Delete(&models.User{}).Error
}
func (d *UserDAO) ListUsers() (users *[]models.User, err error) {
	err = d.db.Find(&users).Error
	return
}

func (d *UserDAO) VerifyUser(userName, password string) (string, error) {
	var user models.User

	// 1. 先根据用户名查找用户
	err := d.
		db.
		Where("name = ?", userName).
		First(&user).
		Error
	if err != nil {
		return "", err
	}

	// 2. 使用 bcrypt 比较密码哈希
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	return user.Role, nil
}
