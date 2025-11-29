package services

import (
	"RedrockHomework/lesson4/database"
	"RedrockHomework/lesson4/models"
	"errors"
)

type StudentService struct{}

func (s *StudentService) CreateStudent(student *models.STUDENT) error {
	result := database.DB.Create(student)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *StudentService) CreateLesson(lesson *models.LESSON) error {
	result := database.DB.Create(lesson)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *StudentService) SelectLesson(studentID int, lessonCode string) error {
	var student models.STUDENT
	if err := database.DB.First(&student, studentID).Error; err != nil {
		return errors.New("学生不存在")
	}
	//寻找学生
	var lesson models.LESSON
	if err := database.DB.First(&lesson, "code = ?", lessonCode).Error; err != nil {
		return errors.New("课程不存在")
	}
	//寻找课程
	var existingRecord models.StudentLesson
	if err := database.DB.Where("student_id = ? AND lesson_id = ?", studentID, lesson.ID).First(&existingRecord).Error; err == nil {
		return errors.New("已经选择过该课程")
	}
	//检查课程是否选满
	var count int64
	database.DB.Where("lesson_id = ?", lesson.ID).Find(&count)
	if int(count) >= lesson.Capacity {
		return errors.New("课程已选满")
	}
	stuentlesson := models.StudentLesson{
		StudentID: studentID,
		LessonID:  lesson.ID,
	}
	return database.DB.Model(&stuentlesson).Create(&stuentlesson).Error
}

func (s *StudentService) DropLesson(studentID int, lessonCode string) error {
	var lesson models.LESSON
	if err := database.DB.First(&lesson, "code = ?", lessonCode).Error; err != nil {
		return errors.New("课程不存在")
	}
	result := database.DB.Where("student_id = ? AND lesson_id = ?", studentID, lesson.ID).Delete(&lesson)
	if result.RowsAffected == 0 {
		return errors.New("选课记录不存在")
	}
	return result.Error
}

func (s *StudentService) GetLessonStudents(lessonCode string) (students []models.STUDENT, err error) {
	var lesson models.LESSON
	if err := database.DB.Where("code = ?", lessonCode).First(&lesson).Error; err != nil {
		err = errors.New("课程不存在")
		return
	}
	var sls []models.StudentLesson
	err = database.DB.Where("lesson_id = ?", lesson.ID).Find(&sls).Error
	if err != nil {
		return
	}
	var student models.STUDENT
	for _, sl := range sls {
		err := database.DB.Where("student_id = ?", sl.StudentID).Find(&student)
		if err != nil {
			continue
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *StudentService) GetStudentLessons(studentID int) (lessons []models.LESSON, err error) {
	var sls []models.StudentLesson
	err = database.DB.Where("lesson_id = ?").Find(&sls).Error
	var lesson models.LESSON
	for _, sl := range sls {
		err := database.DB.Where("student_id = ?", sl.StudentID).Find(&lesson)
		if err != nil {
			continue
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (s *StudentService) FindLesson(lessonCode string) (models.LESSON, error) {
	var lesson models.LESSON
	if err := database.DB.Where("code = ?", lessonCode).First(&lesson).Error; err != nil {
		err = errors.New("课程不存在")
		return lesson, err
	}
	return lesson, nil
}

/*待开发
func (s *StudentService) DeleteStudent(student *models.STUDENT) error {}
func (s *studentService) DeleteLesson(lesson *models.LESSON)error{}
*/
