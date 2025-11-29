package models

import (
	"time"
)

type STUDENT struct {
	Name      string `gorm:"type:varchar(100);not null"`
	ID        int    `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
type LESSON struct {
	Name        string `gorm:"type:varchar(100);not null"`
	ID          int    `gorm:"primaryKey"`
	Code        string `gorm:"type:varchar(20);uniqueIndex;not null"`
	Credit      int    `gorm:"type:int;default:1"`
	Description string `gorm:"type:text"`
	Capacity    int    `gorm:"type:int;default:30"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Duration    string `gorm:"type:text;default:2025/1/1-2025/7/1"` //开始-结束时间
}
type StudentLesson struct {
	StudentID int `gorm:"primaryKey"`
	LessonID  int `gorm:"primaryKey"`
}
