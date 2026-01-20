package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"student_system/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var databaseConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
	var config = databaseConfig
	file, err := os.Open("users.json")
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("无法读取文件: %v", err)
	}

	// 解码 JSON 数据
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("JSON 解码失败: %v", err)
	}

	dsn := config.User + ":" + config.Password + "@" + config.Port + "/" + config.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		// 打印详细的错误信息
		log.Printf("fail to connect with GORM: %v", err)
		fmt.Printf("error type: %T\n", err)
		fmt.Printf("details: %+v\n", err)
		return err
	}

	fmt.Println("Database connected successfully!")

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移
	err = DB.AutoMigrate(
		&models.USER{},
		&models.STUDENT{},
		&models.LESSON{},
		&models.StudentLesson{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	fmt.Println("Database migrated successfully!")
	return nil
}
