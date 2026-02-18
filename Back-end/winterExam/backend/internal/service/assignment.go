package service

import (
	"homeworkSystem/backend/internal/models"
	"homeworkSystem/backend/internal/repository"
	"time"
)

type PublishAssignmentInput struct {
	Title       string
	Description string
	Department  string
	Deadline    time.Time
	AllowLate   bool
	CreatorID   uint64
}

func (s *Service) PublishAssignment(input PublishAssignmentInput) (*models.Assignment, error) {
	if !models.IsValidDepartment(input.Department) {
		return nil, repository.ErrDepartmentNotMatch
	}
	assignment := &models.Assignment{
		Title:       input.Title,
		Description: input.Description,
		Department:  input.Department,
		CreatorID:   input.CreatorID,
		Deadline:    input.Deadline,
		AllowLate:   input.AllowLate,
	}
	if err := s.assignmentRepo.Create(assignment); err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *Service) ListAssignments(department string, page, pageSize int) ([]models.Assignment, int64, error) {
	return s.assignmentRepo.ListByDepartment(department, page, pageSize)
}

func (s *Service) GetAssignmentDetail(id uint64) (*models.Assignment, error) {
	return s.assignmentRepo.FindByID(id)
}

func (s *Service) UpdateAssignment(updater *models.User, id uint64, updates map[string]interface{}) (*models.Assignment, error) {
	// 获取原作业
	assignment, err := s.assignmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 权限检查：必须同部门且为老登
	if updater.Role != models.RoleAdmin {
		return nil, repository.ErrInsufficientPermissions
	}
	if assignment.Department != updater.Department {
		return nil, repository.ErrDepartmentNotMatch
	}

	// 执行更新（带乐观锁，由 repository 处理）
	// 注意：updates 只包含允许修改的字段，由 controller 层过滤
	if err := s.assignmentRepo.Update(assignment); err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *Service) DeleteAssignment(deleter *models.User, id uint64) error {
	assignment, err := s.assignmentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if deleter.Role != models.RoleAdmin {
		return repository.ErrInsufficientPermissions
	}
	if assignment.Department != deleter.Department {
		return repository.ErrDepartmentNotMatch
	}
	return s.assignmentRepo.Delete(id)
}
