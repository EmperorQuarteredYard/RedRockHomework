package repository

import (
	"errors"
	"homeworkSystem/backend/internal/models"

	"gorm.io/gorm"
)

type AssignmentDAO struct {
	db *gorm.DB
}

func NewAssignmentDAO(db *gorm.DB) *AssignmentDAO {
	return &AssignmentDAO{
		db: db,
	}
}

// Publish 给老登发布作业的，不作任何检查
func (d *AssignmentDAO) Publish(ass *models.Assignment) error {
	return d.db.Create(ass).Error
}

// FindByDepartment 按部门筛选，分页查询，不作任何检查
func (d *AssignmentDAO) FindByDepartment(page, pageSize int, department string) (ass *[]models.Assignment, total int64, err error) {
	if err = d.db.Model(&models.Assignment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if offset < 0 || offset > pageSize {
		return nil, 0, err
	}
	err = d.db.Model(ass).Where("department = ?", department).Offset(offset).Limit(pageSize).Find(ass).Error
	if err != nil {
		return nil, 0, err
	}
	return
}

// FindByID 查询单个作业详情，不作任何检查
func (d *AssignmentDAO) FindByID(id uint64) (*models.Assignment, error) {
	result := &models.Assignment{}
	err := d.db.First(result, id).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update 修改作业，检查部门、是否为老资历(已经出现了old light，老登，老资历三个称呼了hh)
func (d *AssignmentDAO) Update(Updater *models.User, newOne *models.Assignment) (err error) {
	var ass *models.Assignment
	if Updater.Role != models.RoleOldLight {
		return errors.New(ErrorInsufficientPermissions)
	}
	ass = &models.Assignment{}
	err = d.db.Where("id = ?", newOne.ID).First(ass).Error
	if err != nil {
		return err
	}
	if ass.Department != Updater.Department {
		return errors.New(ErrorDepartmentNotMatch)
	}
	err = d.db.Model(&models.Assignment{}).Where("id = ?", newOne.ID).Updates(newOne).Error
	return err
}

// Delete 删除作业，检查部门、是否为老资历
func (d *AssignmentDAO) Delete(deleter *models.User, id uint64) (err error) {
	if deleter.Role != models.RoleOldLight {
		return errors.New(ErrorInsufficientPermissions)
	}

	var ass = &models.Assignment{}
	err = d.db.First(ass, id).Error
	if err != nil {
		return err
	}
	if ass.Department != deleter.Department {
		return errors.New(ErrorDepartmentNotMatch)
	}

	return d.db.Delete(ass, id).Error
}
