package controllers

import (
	"errors"
	DAO2 "lesson5/backend/DAO"
	"lesson5/backend/models"

	"gorm.io/gorm"
)

type StudentUser struct {
	SelectionDAO *DAO2.SelectionDAO
	StudentDAO   *DAO2.StudentDAO
	LessonDAO    *DAO2.LessonDAO
}

func NewStudentUser(db *gorm.DB) *StudentUser {
	return &StudentUser{
		SelectionDAO: DAO2.NewSelectionDAO(db),
		StudentDAO:   DAO2.NewStudentDAO(db),
		LessonDAO:    DAO2.NewLessonDAO(db),
	}
}

func (u *StudentUser) Init(db *gorm.DB) {
	u.SelectionDAO = DAO2.NewSelectionDAO(db)
	u.StudentDAO = DAO2.NewStudentDAO(db)
	u.LessonDAO = DAO2.NewLessonDAO(db)
}

func (u *StudentUser) Select(studentID int64, lessonID int64) (bool, error) {
	count, err := u.SelectionDAO.CountSelectionsByLessonID(lessonID)
	if err != nil {
		return false, err
	}

	var exists bool
	exists, err = u.StudentDAO.StudentExist(studentID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("student does not exist")
	}

	exists, err = u.LessonDAO.LessonExist(lessonID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("lesson does not exist")
	}

	capacity, err := u.LessonDAO.GetCapacity(lessonID)
	if err != nil {
		return false, err
	}
	if count < capacity {
		var selection models.Selection
		selection.StudentID = studentID
		selection.LessonID = lessonID
		err = u.SelectionDAO.Create(&selection)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, nil
	}
}

func (u *StudentUser) Drop(studentID int64, lessonID int64) (bool, error) {
	var (
		err       error
		exists    bool
		selection *models.Selection
	)
	exists, err = u.StudentDAO.StudentExist(studentID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("student does not exist")
	}

	exists, err = u.LessonDAO.LessonExist(lessonID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("lesson does not exist")
	}

	selection, err = u.SelectionDAO.GetSelectionByStudentAndLesson(studentID, lessonID)
	if err != nil {
		return false, err
	}
	err = u.SelectionDAO.Delete(selection)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *StudentUser) ListLessons() (*[]models.Lesson, error) {
	return u.LessonDAO.ListLessons()
}

func (u *StudentUser) GetCurrentSelectedLesson(studentID int64) (*[]models.Lesson, error) {
	selections, err := u.SelectionDAO.GetSelectionsByStudentID(studentID)
	var lessons []models.Lesson
	var lesson *models.Lesson
	if err != nil {
		return nil, err
	}
	for _, item := range selections {
		lesson, err = u.LessonDAO.GetLessonByID(item.LessonID)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, *lesson)
	}
	return &lessons, nil
}
