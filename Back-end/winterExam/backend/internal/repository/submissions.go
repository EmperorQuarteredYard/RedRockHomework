package repository

import (
	"fmt"
	"homeworkSystem/backend/internal/models"

	"gorm.io/gorm"
)

type SubmissionDAO struct {
	db *gorm.DB
}

func NewSubmissionDAO(db *gorm.DB) *SubmissionDAO {
	return &SubmissionDAO{
		db: db,
	}
}

func (d *SubmissionDAO) Create(sub *models.Submission) error {
	return d.db.Create(sub).Error
}

func (d *SubmissionDAO) FindBySubmitterID(submitterID uint64) (sub *[]models.Submission, err error) {
	sub = &[]models.Submission{}
	err = d.db.Where("submitter_id = ?", submitterID).Find(sub).Error
	if err != nil {
		return nil, err
	}
	return
}

// FindByDepartment 查看部门提交，不作任何检查
func (d *SubmissionDAO) FindByDepartment(department string) (sub *[]models.Submission, err error) {
	sub = &[]models.Submission{}
	err = d.db.Where("department = ?", department).Find(sub).Error
	if err != nil {
		return nil, err
	}
	return
}

// SetCommentAndScore 批改评语，会检查提交部门是否合法，是不是老登
func (d *SubmissionDAO) SetCommentAndScore(comment string, score int, submissionID uint64, user *models.User) error {
	if user.Role != models.RoleOldLight {
		return fmt.Errorf(ErrorInsufficientPermissions)
	}
	var sub models.Submission
	err := d.db.Where("submitter_id = ?", submissionID).Find(&sub).Error
	if err != nil {
		return err
	}

	if sub.Department != user.Department {
		return fmt.Errorf(ErrorDepartmentNotMatch)
	}

	sub.Comment = comment
	sub.Score = score
	return d.db.Updates(&sub).Error
}

// SetExcellent 标记优秀，会检查提交部门是否合法，是不是老登
func (d *SubmissionDAO) SetExcellent(submissionID uint64, user *models.User) error {
	if user.Role != models.RoleOldLight {
		return fmt.Errorf(ErrorInsufficientPermissions)
	}
	var sub models.Submission
	err := d.db.Where("submitter_id = ?", submissionID).Find(&sub).Error
	if err != nil {
		return err
	}

	if sub.Department != user.Department {
		return fmt.Errorf(ErrorDepartmentNotMatch)
	}
	if sub.Excellent {
		return nil
	}
	return d.db.Model(&models.Submission{}).Where("id = ?", submissionID).Update("excellent", true).Error

}

// GetExcellent 优秀作业展示，不作任何检查
func (d *SubmissionDAO) GetExcellent(department string) (sub *[]models.Submission, err error) {
	sub = &[]models.Submission{}
	err = d.db.Model(&models.Submission{}).Where("excellent = ? and department = ?", true, department).Find(sub).Error
	if err != nil {
		return nil, err
	}
	return
}
