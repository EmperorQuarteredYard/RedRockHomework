package configs

import (
	"encoding/json"
	"homeworkSystem/backend/internal/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DB  DSNConfig `json:"db"`
	JWT JWTConfig `json:"jwt"`
}

type DSNConfig struct {
	User     string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"database"`
}

type JWTConfig struct {
	AccessSecret  string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
	Issuer        string `json:"issuer"`
	AccessExpire  int    `json:"access_expire"`
	RefreshExpire int    `json:"refresh_expire"`
}

func LoadConfig() *Config {
	// 实际项目中应从环境变量或配置文件读取
	var config Config
	file, err := os.Open("backend/configs/config.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return nil
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return nil
	}
	return &config
	//return &Config{
	//	DB: DSNConfig{
	//		User:     "root",
	//		Password: "123456",
	//		Host:     "127.0.0.1",
	//		Port:     "3306",
	//		DBName:   "homework",
	//	},
	//	JWT: JWTConfig{
	//		AccessSecret:  "access_secret_key",
	//		RefreshSecret: "refresh_secret_key",
	//		AccessExpire:  2, // 2小时
	//		RefreshExpire: 7, // 7天
	//	},
	//}
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
