package repository

import (
	"errors"
	"homeworkSystem/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssignmentRepo struct {
	db *gorm.DB
}

func NewAssignmentRepo(db *gorm.DB) *AssignmentRepo {
	return &AssignmentRepo{db: db}
}

func (r *AssignmentRepo) Create(assignment *models.Assignment) error {
	return r.db.Create(assignment).Error
}

func (r *AssignmentRepo) FindByID(id uint64) (*models.Assignment, error) {
	var assignment models.Assignment
	err := r.db.First(&assignment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &assignment, nil
}

func (r *AssignmentRepo) ListByDepartment(department string, page, pageSize int) ([]models.Assignment, int64, error) {
	var list []models.Assignment
	var total int64
	db := r.db.Model(&models.Assignment{})
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

// Update 带乐观锁并发控制
func (r *AssignmentRepo) Update(assignment *models.Assignment) error {
	return r.db.Model(assignment).Clauses(clause.Returning{}).Updates(assignment).Error
}

func (r *AssignmentRepo) Delete(id uint64) error {
	return r.db.Delete(&models.Assignment{}, id).Error
}
