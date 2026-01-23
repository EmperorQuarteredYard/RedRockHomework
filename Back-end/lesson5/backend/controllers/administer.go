package controllers

import (
	"lesson5/backend/DAO"
	"lesson5/backend/models"

	"gorm.io/gorm"
)

type AdministerUser struct {
	UserDAO *DAO.UserDAO
	*StudentUser
}

func NewAdministerUser(db *gorm.DB) *AdministerUser {
	userDao := DAO.NewUserDAO(db)
	return &AdministerUser{UserDAO: userDao, StudentUser: NewStudentUser(db)}
}

func (u *AdministerUser) CreateStudent(student *models.Student) error {
	return u.StudentDAO.CreateStudent(student)
}

func (u *AdministerUser) UpdateStudent(student *models.Student) error {
	return u.StudentDAO.UpdateStudent(student)
}

func (u *AdministerUser) ListStudent() (students *[]models.Student, err error) {
	return u.StudentDAO.ListStudents()
}

func (u *AdministerUser) DeleteStudent(studentID int64) error {
	return u.StudentDAO.DeleteStudentByID(studentID)
}

func (u *AdministerUser) CreateLesson(lesson *models.Lesson) error {
	return u.LessonDAO.CreateLesson(lesson)
}

func (u *AdministerUser) UpdateLesson(lesson *models.Lesson) error {
	return u.LessonDAO.UpdateLesson(lesson)
}

func (u *AdministerUser) DeleteLesson(lessonID int64) error {
	return u.LessonDAO.DeleteLesson(lessonID)
}

func (u *AdministerUser) Login(userName, password string) (int64, error) {
	role, err := u.UserDAO.VerifyUser(userName, password)
	if err != nil {
		return 0, err
	}
	if role == "student" {
		return 1, nil
	}
	if role == "teacher" {
		return 2, nil
	}
	if role == "administrator" {
		return 4, nil
	}
	return 0, nil
}

func (u *AdministerUser) ListSelection() (selections *[]models.Selection, err error) {
	selections, err = u.SelectionDAO.ListSelections()
	return
}
