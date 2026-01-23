package models

import (
	"time"
)

type AllModels interface {
	STUDENT |
		LESSON |
		StudentLesson |
		USER
} //有了接口之后就好啦
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
	Duration    string `gorm:"type:text"` //开始-结束时间
}
type StudentLesson struct {
	StudentID int `gorm:"primaryKey"`
	LessonID  int `gorm:"primaryKey"`
}
