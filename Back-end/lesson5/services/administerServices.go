package services

import (
	"student_system/database"
	"student_system/models"
)

type AdministerServices struct{}

func (a *AdministerServices) CreateStudent(student *models.STUDENT) error {
	result := database.DB.Create(student)
	return result.Error
}

func (a *AdministerServices) CreateStudents(students *[]models.STUDENT) error {
	for _, student := range *students {
		err := a.CreateStudent(&student)
		if err != nil {
			return err
		}
	}
	return nil
}


func (a *AdministerServices) FindStudent(studentName string) (student *models.STUDENT, err error) {
	err = database.DB.Where("name = ?", studentName).Find(&student).Error
	return
}

func (a *AdministerServices) GetStudents() (students []models.STUDENT, err error) {
	err = database.DB.Table("students").Find(&students).Error
	return
}

func (a *AdministerServices) CreateSLesson(lesson models.LESSON) error {
	result := database.DB.Create(&lesson)
	return result.Error
}

func (a *AdministerServices) CreateSLessons(lessons []models.LESSON) error {
	for _, lesson := range lessons {
		err := a.CreateSLesson(lesson)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AdministerServices) GetStudentSelections() (studentSelections map[models.STUDENT][]models.LESSON, err error) {
	var studentLessons []models.StudentLesson
	err = database.DB.Table("student_lesson").Find(&studentLessons).Error
	if err != nil {
		return
	}
	var lesson models.LESSON
	var student models.STUDENT
	for _, studentLesson := range studentLessons {
		err = database.DB.Table("students").Where("student_id = ?", studentLesson.StudentID).Find(&student).Error
		if err != nil {
			return
		}
		err = database.DB.Table("lessons").Where("lesson_id = ?", studentLesson.StudentID).Find(&lesson).Error
		if err != nil {
			return
		}
		if studentSelections[student] == nil {
			studentSelections[student] = []models.LESSON{}
		}
		studentSelections[student] = append(studentSelections[student], lesson)
	}
	return
}

func (a *AdministerServices) GetLessonStudents(lesson models.LESSON) (students []models.STUDENT, err error) {
	var sls []models.StudentLesson
	err = database.DB.Where("lesson_id = ?", lesson.ID).Find(&sls).Error
	if err != nil {
		return
	}
	var student models.STUDENT
	for _, sl := range sls {
		err = database.DB.Where("student_id = ?", sl.StudentID).Find(&student).Error
		if err != nil {
			panic(err)
		}
		students = append(students, student)
	}
	return students, nil
}

func (a *AdministerServices) GetStudentLessons(student models.StudentLesson) (lessons []models.LESSON, err error) {
	var sls []models.StudentLesson
	err = database.DB.Where("lesson_id = ?").Find(&sls).Error
	var lesson models.LESSON
	for _, sl := range sls {
		err = database.DB.Where("student_id = ?", sl.StudentID).Find(&lesson).Error
		if err != nil {
			panic(err)
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (a *AdministerServices) DeleteStudent(student models.STUDENT) error {
	result := database.DB.Delete(&student)
	return result.Error
}

func (a *AdministerServices) DeleteStudents(students []models.STUDENT) error {
	for _, student := range students {
		err := a.DeleteStudent(student)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AdministerServices) DeleteLesson(lesson models.LESSON) error {
	result := database.DB.Delete(&lesson)
	return result.Error
}

func (a *AdministerServices) DeleteLessons(lessons []models.LESSON) error {
	for _, lesson := range lessons {
		err := a.DeleteLesson(lesson)
		if err != nil {
			return err
		}
	}
	return nil
}
