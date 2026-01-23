package models

import (
	"time"
)

type Student struct {
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	ID        int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}
type Lesson struct {
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	ID          int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Credit      int64     `gorm:"type:bigint;default:1" json:"credit"`
	Description string    `gorm:"type:text" json:"description"`
	Capacity    int64     `gorm:"type:bigint;default:30" json:"capacity"`
	CreatedAt   time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null" json:"updated_at"`
	Duration    string    `gorm:"type:text"` //开始-结束时间
}
type Selection struct {
	ID        int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	StudentID int64     `gorm:"type:bigint;primaryKey" json:"student_id"`
	LessonID  int64     `gorm:"type:bigint;primaryKey" json:"lesson_id"`
	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}
type User struct {
	Name         string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	ID           int64     `gorm:"type:bigint;primaryKey;autoIncrement" json:"id"`
	Role         string    `gorm:"type:varchar(100);not null" json:"role"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}
