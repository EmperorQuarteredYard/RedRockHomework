package database

import (
	"encoding/json"
	"fmt"
	"lesson5/backend/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type databaseConfig struct {
	USERNAME string `json:"username"`
	PASSWORD string `json:"password"`
	DATABASE string `json:"database"`
	PORT     string `json:"port"`
	HOST     string `json:"host"`
}

func getConfig() (*databaseConfig, error) {
	file, err := os.Open("backend/config/database.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var newConfig databaseConfig
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&newConfig)
	if err != nil {
		return nil, err
	}
	return &newConfig, nil
}

func InitDB() *gorm.DB {
	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
	}
	if config == nil {
		fmt.Println("config is nil")
		return nil
	}
	dsn := config.USERNAME + ":" + config.PASSWORD + "@tcp(" + config.HOST + ":" + config.PORT + ")/" + config.DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("fail to connect to database\nreason:", err)
		return nil
	} else {
		fmt.Println("successfully connect to database")
	}

	err = DB.AutoMigrate(&models.Student{}, &models.Lesson{}, &models.Selection{}, &models.User{})

	if err != nil {
		fmt.Println("Fail to migrate model\nreason:", err)
	} else {
		fmt.Println("successfully migrate model")
	}
	return DB
}
