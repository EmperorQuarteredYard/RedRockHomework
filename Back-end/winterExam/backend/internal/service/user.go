package service

import (
	"errors"
	"homeworkSystem/backend/internal/models"
	"homeworkSystem/backend/internal/repository"
	jwt "homeworkSystem/backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(username, password, nickname, department string) (*models.User, error) {
	// 校验部门
	if !models.IsValidDepartment(department) {
		return nil, repository.ErrDepartmentNotMatch
	}

	// 检查用户名是否已存在
	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return nil, repository.ErrDuplicateEntry
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashed),
		Nickname:     nickname,
		Role:         models.RoleStudent,
		Department:   department,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(username, password string) (*models.User, string, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", "", err
	}
	if err := s.userRepo.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, "", "", err
	}
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.Department, user.Role)
	if err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}

func (s *Service) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	// 重新生成双 token
	return jwt.GenerateToken(claims.UserID, claims.Department, claims.Role)
}

func (s *Service) GetUserByID(id uint64) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *Service) DeleteAccount(userID uint64, password string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if err := s.userRepo.ComparePassword(user.PasswordHash, password); err != nil {
		return err
	}
	return s.userRepo.SoftDelete(userID)
}

func (s *Service) GetNicknameByID(userID uint64) (string, error) {
	return s.userRepo.GetNicknameByID(userID)
} // PromoteToAdmin 提升用户为管理员
func (s *Service) PromoteToAdmin(userID uint64) error {
	return s.userRepo.PromoteToAdmin(userID)
}
