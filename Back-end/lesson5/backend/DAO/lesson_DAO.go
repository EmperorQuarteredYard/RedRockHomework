package DAO

import (
	"errors"
	"fmt"
	"lesson5/backend/models"

	"gorm.io/gorm"
)

type LessonDAO struct {
	*BaseDAO
}

func NewLessonDAO(db *gorm.DB) *LessonDAO {
	if db == nil {
		fmt.Println("NewLessonDAO Called: db不能为空")
		panic(errors.New("NewLessonDAO Called: db不能为空"))
		return nil
	}
	return &LessonDAO{NewBaseDAO(db)}
}

func (d *LessonDAO) GetAllLessons(lesson *[]models.Lesson) error {
	return d.db.Find(lesson).Error
}

func (d *LessonDAO) ListLessons() (lessons *[]models.Lesson, err error) {
	err = d.db.Find(&lessons).Error
	return
}

func (d *LessonDAO) LessonExist(id int64) (bool, error) {
	return d.Exist(&models.Lesson{}, id)
}

func (d *LessonDAO) GetLessonByID(id int64) (*models.Lesson, error) {
	var lesson models.Lesson
	err := d.db.Where("id = ?", id).First(&lesson).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (d *LessonDAO) CreateLesson(lesson *models.Lesson) error {
	return d.Create(lesson)
}

func (d *LessonDAO) UpdateLesson(lesson *models.Lesson) error {
	return d.Update(lesson)
}

func (d *LessonDAO) DeleteLesson(id int64) error {
	return d.db.Delete(&models.Lesson{}, id).Error
}

func (d *LessonDAO) GetCapacity(lessonID int64) (int64, error) {
	lesson, err := d.GetLessonByID(lessonID)
	if err != nil {
		return 0, err
	}
	return lesson.Capacity, nil
}
