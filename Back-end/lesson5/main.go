package main

import (
	"fmt"
	"student_system/database"
	"student_system/models"
	"student_system/services"
)

func main() {
	database.InitDB()
	var s services.AdministerServices
	var student = models.STUDENT{
		Name: "张三",
	}
	database.DB.AutoMigrate(&student)
	err := s.CreateStudent(&student)
	if err != nil {
		fmt.Println("error:", err)
	}
	students, _ := s.GetStudents()
	for _, stu := range students {
		fmt.Println(stu)
	}
	students, _ = s.GetStudents()
	for _, stu := range students {
		fmt.Println(stu)
	}
}
