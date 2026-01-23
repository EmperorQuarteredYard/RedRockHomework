package DAO

import (
	"errors"
	"fmt"
	"lesson5/backend/models"

	"gorm.io/gorm"
)

type SelectionDAO struct {
	*BaseDAO
}

func NewSelectionDAO(db *gorm.DB) *SelectionDAO {
	if db == nil {
		fmt.Println("NewSelectionDAO Called: db不能为空")
		panic(errors.New("NewSelectionDAO Called: db不能为空"))
		return nil
	}
	return &SelectionDAO{NewBaseDAO(db)}
}

func (d *SelectionDAO) ListSelections() (*[]models.Selection, error) {
	var selections []models.Selection
	err := d.db.Find(&selections).Error
	return &selections, err
}

func (d *SelectionDAO) GetSelectionByID(id int64) (*models.Selection, error) {
	var selection models.Selection
	err := d.db.Where("id = ?", id).First(&selection).Error
	if err != nil {
		return nil, err
	}
	return &selection, nil
}

func (d *SelectionDAO) GetSelectionByStudentAndLesson(studentID, lessonID int64) (*models.Selection, error) {
	var selection models.Selection
	err := d.db.Where("student_id = ? and lesson_id = ?", studentID, lessonID).First(&selection).Error
	if err != nil {
		return nil, err
	}
	return &selection, nil
}

func (d *SelectionDAO) CreateSelection(selection *models.Selection) error {
	return d.db.Create(selection).Error
}

func (d *SelectionDAO) UpdateSelection(selection *models.Selection) error {
	return d.db.Save(selection).Error
}

func (d *SelectionDAO) DeleteSelection(id int64) error {
	return d.db.Delete(&models.Selection{}, id).Error
}

func (d *SelectionDAO) GetSelectionsByStudentID(studentID int64) ([]models.Selection, error) {
	var selections []models.Selection
	err := d.db.Where("student_id = ?", studentID).Find(&selections).Error
	return selections, err
}

func (d *SelectionDAO) GetSelectionsByLessonID(lessonID int64) ([]models.Selection, error) {
	var selections []models.Selection
	err := d.db.Where("lesson_id = ?", lessonID).Find(&selections).Error
	return selections, err
}

func (d *SelectionDAO) CountSelectionsByStudentID(studentID int64) (int64, error) {
	var count int64
	err := d.db.Model(&models.Selection{}).Where("student_id = ?", studentID).Count(&count).Error
	return count, err
}

func (d *SelectionDAO) CountSelectionsByLessonID(lessonID int64) (int64, error) {
	var count int64
	err := d.db.Model(&models.Selection{}).Where("lesson_id = ?", lessonID).Count(&count).Error
	return count, err
}
