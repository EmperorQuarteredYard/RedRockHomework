package service

import (
	"homeworkSystem/backend/internal/repository"

	"gorm.io/gorm"
)

type Service struct {
	assignmentDao *repository.AssignmentDAO
	userDao       *repository.UserDao
	submissionDao *repository.SubmissionDAO
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		assignmentDao: repository.NewAssignmentDAO(db),
		userDao:       repository.NewUserDao(db),
		submissionDao: repository.NewSubmissionDAO(db),
	}
}
