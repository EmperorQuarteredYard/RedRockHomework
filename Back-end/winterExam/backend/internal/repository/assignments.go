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

// ListByDepartment 当department为空时，返回所有记录；否则，返回department对应的记录
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

// Update 使用事务锁处理并发
func (r *AssignmentRepo) Update(assignment *models.Assignment) error {
	var assignments models.Assignment
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&assignments).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Assignment{}).Updates(assignment).Error; err != nil {
			return err
		}
		if err := tx.First(assignment, assignment.ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateByMap 使用事务锁更新作业字段
func (r *AssignmentRepo) UpdateByMap(id uint64, updates map[string]interface{}) (*models.Assignment, error) {
	var assignment models.Assignment
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&assignment, id).Error; err != nil {
			return err
		}
		if err := tx.Model(&assignment).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.First(&assignment, id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *AssignmentRepo) Delete(id uint64) error {
	return r.db.Delete(&models.Assignment{}, id).Error
}
