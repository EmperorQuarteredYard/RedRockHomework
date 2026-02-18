package repository

import (
	"errors"
	"homeworkSystem/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *SubmissionRepo) FindByStudentHomework(studentID, homeworkID uint64) ([]models.Submission, error) {
	var list []models.Submission
	err := r.db.Model(&models.Submission{}).Where("student_id = ? AND homework_id = ?", studentID, homeworkID).Order("created_at desc").Find(&list).Error
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

func (r *SubmissionRepo) FindByDepartmentHomework(department string, homeworkID uint64) ([]models.Submission, error) {
	var list []models.Submission
	err := r.db.Where("department = ? AND homework_id = ?", department, homeworkID).Find(&list).Error
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

func (r *SubmissionRepo) ListExcellentByDepartment(department string, page, pageSize int) ([]models.Submission, int64, error) {
	var list []models.Submission
	var total int64
	db := r.db.Model(&models.Submission{})
	if department != "" {
		db = db.Where("department = ?", department)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&list).Error
	return list, total, err
}
func (r *SubmissionRepo) UpdateByMap(id uint64, updates map[string]interface{}) (*models.Submission, error) {
	var submission models.Submission
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 锁定行
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&submission, id).Error; err != nil {
			return err
		}
		// 执行更新
		if err := tx.Model(&submission).Updates(updates).Error; err != nil {
			return err
		}
		// 重新查询以获取更新后的值
		if err := tx.First(&submission, id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *SubmissionRepo) CountByStudentHomework(studentID, homeworkID uint64) (count int64, err error) {
	err = r.db.Model(&models.Submission{}).Where("student_id = ? and homework_id = ?", studentID, homeworkID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return
}

func (r *SubmissionRepo) CountByHomework(homeworkID uint64) (count int64, err error) {
	err = r.db.Model(&models.Submission{}).Where("homework_id = ?", homeworkID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return
}
