package repository

import (
	"homeworkSystem/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (d *UserDao) Register(userName, password, department, nickname, role string) (user *models.User, err error) {
	_, err = models.GetDepartment(department)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user = &models.User{
		Name:         userName,
		PasswordHash: string(hash),
		Role:         role,
		Department:   department,
		Nickname:     nickname,
	}
	err = d.db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *UserDao) Login(userName, password string) (user *models.User, err error) {
	user = &models.User{}
	err = d.db.Where("name = ?", userName).First(user).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *UserDao) SoftDelete(userID uint64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{}
	err = d.db.Where("id = ?", userID).First(user).Error
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return err
	}

	return d.db.Where("id = ? and password_hash = ?", userID, hash).Delete(&models.User{}).Error
}

func (d *UserDao) FindByID(ID uint64) (*models.User, error) {
	user := &models.User{}
	result := d.db.Where("id = ?", ID).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
