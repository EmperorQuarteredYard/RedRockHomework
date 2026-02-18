package service

import (
	"gorm.io/gorm"
	"homeworkSystem/backend/internal/repository"
)

type Service struct {
	userRepo       *repository.UserRepo
	assignmentRepo *repository.AssignmentRepo
	submissionRepo *repository.SubmissionRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		userRepo:       repository.NewUserRepo(db),
		assignmentRepo: repository.NewAssignmentRepo(db),
		submissionRepo: repository.NewSubmissionRepo(db),
	}
}
