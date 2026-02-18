package service

import (
	"homeworkSystem/backend/internal/models"
	"homeworkSystem/backend/internal/repository"
	"time"
)

type SubmitInput struct {
	HomeworkID uint64
	StudentID  uint64
	Content    string
	FileURL    string
}

func (s *Service) Submit(input SubmitInput) (*models.Submission, error) {
	// 检查作业是否存在
	assignment, err := s.assignmentRepo.FindByID(input.HomeworkID)
	if err != nil {
		return nil, err
	}

	// 检查学生部门与作业部门是否匹配
	student, err := s.userRepo.FindByID(input.StudentID)
	if err != nil {
		return nil, err
	}
	if student.Department != assignment.Department {
		return nil, repository.ErrDepartmentNotMatch
	}

	// 检查截止时间
	now := time.Now()
	isLate := now.After(assignment.Deadline)
	if isLate && !assignment.AllowLate {
		return nil, repository.ErrSubmitLate // 需定义新错误
	}

	submission := &models.Submission{
		HomeworkID:  input.HomeworkID,
		StudentID:   input.StudentID,
		Content:     input.Content,
		FileURL:     input.FileURL,
		IsLate:      isLate,
		SubmittedAt: now,
	}
	if err := s.submissionRepo.Create(submission); err != nil {
		return nil, err
	}
	return submission, nil
}

func (s *Service) GetMySubmissions(studentID uint64) ([]models.Submission, error) {
	return s.submissionRepo.FindByStudent(studentID)
}

func (s *Service) GetHomeworkSubmissions(homeworkID uint64) ([]models.Submission, error) {
	return s.submissionRepo.FindByHomework(homeworkID)
}

func (s *Service) ReviewSubmission(reviewer *models.User, submissionID uint64, score int, comment string, isExcellent bool) (*models.Submission, error) {
	submission, err := s.submissionRepo.FindByID(submissionID)
	if err != nil {
		return nil, err
	}

	// 检查权限：必须同部门老登
	if reviewer.Role != models.RoleAdmin {
		return nil, repository.ErrInsufficientPermissions
	}
	// 需要通过 homework 获取部门
	homework, err := s.assignmentRepo.FindByID(submission.HomeworkID)
	if err != nil {
		return nil, err
	}
	if homework.Department != reviewer.Department {
		return nil, repository.ErrDepartmentNotMatch
	}

	submission.Score = score
	submission.Comment = comment
	submission.IsExcellent = isExcellent
	submission.ReviewerID = reviewer.ID
	submission.ReviewedAt = time.Now()

	if err := s.submissionRepo.Update(submission); err != nil {
		return nil, err
	}
	return submission, nil
}

func (s *Service) MarkExcellent(reviewer *models.User, submissionID uint64, isExcellent bool) (*models.Submission, error) {
	submission, err := s.submissionRepo.FindByID(submissionID)
	if err != nil {
		return nil, err
	}

	homework, err := s.assignmentRepo.FindByID(submission.HomeworkID)
	if err != nil {
		return nil, err
	}
	if reviewer.Role != models.RoleAdmin {
		return nil, repository.ErrInsufficientPermissions
	}
	if homework.Department != reviewer.Department {
		return nil, repository.ErrDepartmentNotMatch
	}

	submission.IsExcellent = isExcellent
	submission.ReviewerID = reviewer.ID
	submission.ReviewedAt = time.Now()

	if err := s.submissionRepo.Update(submission); err != nil {
		return nil, err
	}
	return submission, nil
}

func (s *Service) GetExcellentSubmissions(department string) ([]models.Submission, error) {
	return s.submissionRepo.FindExcellent(department)
}
