package database

import (
	"fmt"
	"lesson4/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var databaseConfig struct {
	USERNAME string
	PASSWORD string
	DATABASE string
	PORT     string
	HOST     string
}
config :=databaseConfig{}
func getConfig() {
	file,err:=os.Open("config/database.json")
}
func InitDB() {
	dsn := "RedRockHomework_ClassSelection:BestRedRock@tcp(your_host:8080)/student_system?charset=utf8mb4&parseTime=True&loc=Local"
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
