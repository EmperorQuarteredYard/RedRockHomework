package DAO

import (
	"errors"
	"fmt"
	"lesson5/backend/models"

	"gorm.io/gorm"
)

type StudentDAO struct {
	*BaseDAO
}

func NewStudentDAO(db *gorm.DB) *StudentDAO {
	if db == nil {
		fmt.Println("NewStudentDAO Called: db不能为空")
		panic(errors.New("NewStudentDAO Called: db不能为空"))
		return nil
	}
	return &StudentDAO{NewBaseDAO(db)}
}

func (d *StudentDAO) DeleteStudentByID(studentID int64) error {
	return d.db.Delete(&models.Student{}, studentID).Error
}

func (d *StudentDAO) ListStudents() (*[]models.Student, error) {
	var students []models.Student
	err := d.db.Find(&students).Error
	if err != nil {
		return nil, err
	}
	return &students, nil
}

func (d *StudentDAO) StudentExist(id int64) (bool, error) {
	return d.Exist(&models.Student{}, id)
}

func (d *StudentDAO) GetStudentByID(id int64) (student *models.Student, err error) {
	err = d.db.Where("id = ?", id).First(student).Error
	return
}

func (d *StudentDAO) GetStudentByName(name string) (student *models.Student, err error) {
	err = d.db.Where("name = ?", name).First(student).Error
	return
}

func (d *StudentDAO) CreateStudent(student *models.Student) error {
	return d.Create(student)
}

func (d *StudentDAO) UpdateStudent(student *models.Student) error {
	return d.Update(student)
}

func (d *StudentDAO) DeleteStudent(student *models.Student) error {

	return d.Delete(student)
}
