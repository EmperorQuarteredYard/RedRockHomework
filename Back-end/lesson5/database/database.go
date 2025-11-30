package database

import (
	"fmt"
	"student_system/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "RedRockHomework_ClassSelection:BestRedRock@tcp(your_host:8080)/student_system?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
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
