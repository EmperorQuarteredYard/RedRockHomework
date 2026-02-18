package repository

import (
	"errors"
	"homeworkSystem/backend/internal/models"

	"gorm.io/gorm"
)

type SubmissionRepo struct {
	db *gorm.DB
}

func NewSubmissionRepo(db *gorm.DB) *SubmissionRepo {
	return &SubmissionRepo{db: db}
}

func (r *SubmissionRepo) Create(submission *models.Submission) error {
	return r.db.Create(submission).Error
}

func (r *SubmissionRepo) FindByID(id uint64) (*models.Submission, error) {
	var sub models.Submission
	err := r.db.First(&sub, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &sub, nil
}

func (r *SubmissionRepo) FindByStudent(studentID uint64) ([]models.Submission, error) {
	var list []models.Submission
	err := r.db.Where("student_id = ?", studentID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *SubmissionRepo) FindByHomework(homeworkID uint64) ([]models.Submission, error) {
	var list []models.Submission
	err := r.db.Where("homework_id = ?", homeworkID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *SubmissionRepo) FindByDepartment(department string) ([]models.Submission, error) {
	var list []models.Submission
	err := r.db.Where("department = ?", department).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *SubmissionRepo) Update(submission *models.Submission) error {
	return r.db.Save(submission).Error
}

func (r *SubmissionRepo) FindExcellent(department string) ([]models.Submission, error) {
	var list []models.Submission
	db := r.db.Where("is_excellent = ?", true)
	if department != "" {
		db = db.Where("department = ?", department)
	}
	err := db.Order("score desc").Find(&list).Error
	return list, err
}
