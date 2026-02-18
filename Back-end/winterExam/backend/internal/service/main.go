// Package service 业务的service层，实际上也封装了controller层的内容
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
