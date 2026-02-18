package configs

import (
	"homeworkSystem/backend/internal/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DB  DSNConfig
	JWT JWTConfig
}

type DSNConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpire  int // 小时
	RefreshExpire int // 天
}

func LoadConfig() *Config {
	// 实际项目中应从环境变量或配置文件读取
	return &Config{
		DB: DSNConfig{
			User:     "root",
			Password: "123456",
			Host:     "127.0.0.1",
			Port:     "3306",
			DBName:   "homework",
		},
		JWT: JWTConfig{
			AccessSecret:  "access_secret_key",
			RefreshSecret: "refresh_secret_key",
			AccessExpire:  2, // 2小时
			RefreshExpire: 7, // 7天
		},
	}
}

func InitDB(cfg *Config) *gorm.DB {
	dsn := cfg.DB.User + ":" + cfg.DB.Password + "@tcp(" + cfg.DB.Host + ":" + cfg.DB.Port + ")/" + cfg.DB.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&models.User{},
		&models.Assignment{},
		&models.Submission{},
	); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	return db
}
