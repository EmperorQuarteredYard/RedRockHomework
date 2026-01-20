package models

import (
	"time"
)

type STUDENT struct {
	Name      string    `gorm:"type:varchar(100);not null"`
	ID        int       `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}
type LESSON struct {
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	ID          int       `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Credit      int       `gorm:"type:int;default:1" json:"credit"`
	Description string    `gorm:"type:text" json:"description"`
	Capacity    int       `gorm:"type:int;default:30" json:"capacity"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null" json:"updated_at"`
	Duration    string    `gorm:"type:text;default:2025/1/1-2025/7/1"` //开始-结束时间
}
type StudentLesson struct {
	StudentID int       `gorm:"primaryKey" json:"student_id"`
	LessonID  int       `gorm:"primaryKey" json:"lesson_id"`
	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}
