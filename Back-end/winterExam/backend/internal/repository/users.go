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

func (d *UserDao) Register(userName, password, role, department string) (user *models.User, err error) {
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

func (d *UserDao) SoftDelete(userID string) error {
	result := d.db.Delete(&models.User{}, userID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *UserDao) FindByID(ID uint64) (*models.User, error) {
	user := &models.User{}
	result := d.db.Where("id = ?", ID).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
