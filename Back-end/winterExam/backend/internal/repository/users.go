package repository

import (
	"errors"
	"homeworkSystem/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) SoftDelete(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}

// ComparePassword 验证密码
func (r *UserRepo) ComparePassword(hashedPwd, plainPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
}

func (r *UserRepo) GetNicknameByID(id uint64) (string, error) {
	var nickname string
	err := r.db.Model(&models.User{}).Where("id = ?", id).Select("nickname").Scan(&nickname).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return nickname, nil
}
func (r *UserRepo) PromoteToAdmin(userID uint64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", models.RoleOldLight).Error
}

// Update 更新用户信息
func (r *UserRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}
