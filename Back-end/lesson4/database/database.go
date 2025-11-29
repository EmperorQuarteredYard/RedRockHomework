package database

import (
	"RedrockHomework/lesson4/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:123456@tcp(your_host:8080)/student_system?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("fail to connect to database\nreason:", err)
		return
	} else {
		fmt.Println("successfully connect to database")
	}

	err = DB.AutoMigrate(&models.STUDENT{}, &models.LESSON{}, &models.StudentLesson{})

	if err != nil {
		fmt.Println("Fail to migrate model\nreason:", err)
	} else {
		fmt.Println("successfully migrate model")
	}
}
