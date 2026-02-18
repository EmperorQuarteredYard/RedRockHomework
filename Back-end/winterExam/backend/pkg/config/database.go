package config

import (
	"encoding/json"
	"fmt"
	"homeworkSystem/backend/internal/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	initialized bool = false
	db          *gorm.DB
)

type databaseConfig struct {
	USERNAME string
	PASSWORD string
	DATABASE string
	PORT     string
	HOST     string
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

func GetDatabase() *gorm.DB {
	if !initialized {
		config, err := getConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if config == nil {
			fmt.Println("config is nil")
			return nil
		}
		dsn := config.USERNAME + ":" + config.PASSWORD + "@tcp(" + config.HOST + ":" + config.PORT + ")/" + config.DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if err != nil {
			fmt.Println("fail to connect to database\nreason:", err)
			return nil
		} else {
			fmt.Println("successfully connect to database")
		}

		err = db.AutoMigrate(&models.User{}, &models.Submission{}, &models.Assignment{})

		if err != nil {
			fmt.Println("Fail to migrate model\nreason:", err)
		} else {
			fmt.Println("successfully migrate model")
		}
		initialized = true
	}
	return db
}
