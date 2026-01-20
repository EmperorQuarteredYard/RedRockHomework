package services

import (
	"errors"
	"student_system/database"
	"student_system/models"
)

type StudentService struct {
	student models.STUDENT
}

func (s *StudentService) SelectLesson(lessonCode string) error {
	var lesson models.LESSON
	if err := database.DB.First(&lesson, "code = ?", lessonCode).Error; err != nil {
		return errors.New("课程不存在")
	}
	//寻找课程
	var existingRecord models.StudentLesson
	if err := database.DB.Where("student_id = ? AND lesson_id = ?", s.student.ID, lesson.ID).First(&existingRecord).Error; err == nil {
		return errors.New("已经选择过该课程")
	}
	//检查课程是否选满
	var count int64
	database.DB.Where("lesson_id = ?", lesson.ID).Find(&count)
	if int(count) >= lesson.Capacity {
		return errors.New("课程已选满")
	}
	stuentlesson := models.StudentLesson{
		StudentID: s.student.ID,
		LessonID:  lesson.ID,
	}
	return database.DB.Model(&stuentlesson).Create(&stuentlesson).Error
}

func (s *StudentService) DropLesson(lessonCode string) error {
	var lesson models.LESSON
	if err := database.DB.First(&lesson, "code = ?", lessonCode).Error; err != nil {
		return errors.New("课程不存在")
	}
	result := database.DB.Where("student_id = ? AND lesson_id = ?", s.student.ID, lesson.ID).Delete(&lesson)
	if result.RowsAffected == 0 {
		return errors.New("选课记录不存在")
	}
	return result.Error
}

func (s *StudentService) GetStudentLessons() (lessons []models.LESSON, err error) {
	var sls []models.StudentLesson
	err = database.DB.Where("student_id = ?", s.student.ID).Find(&sls).Error
	var lesson models.LESSON
	for _, sl := range sls {
		err = database.DB.Where("lesson_id = ?", sl.LessonID).Find(&lesson).Error
		if err != nil {
			panic(err)
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}
