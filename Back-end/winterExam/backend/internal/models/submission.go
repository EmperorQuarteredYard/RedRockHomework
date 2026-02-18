package models

import (
	"gorm.io/gorm"
	"time"
)

type Submission struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	HomeworkID  uint64         `json:"homework_id" gorm:"not null"`
	StudentID   uint64         `json:"student_id" gorm:"not null"`
	Content     string         `json:"content" gorm:"type:text"`
	FileURL     string         `json:"file_url"`
	IsLate      bool           `json:"is_late" gorm:"default:false"`
	Score       int            `json:"score"`
	Comment     string         `json:"comment" gorm:"type:text"`
	IsExcellent bool           `json:"is_excellent" gorm:"default:false"`
	ReviewerID  uint64         `json:"reviewer_id"`
	SubmittedAt time.Time      `json:"submitted_at"`
	ReviewedAt  time.Time      `json:"reviewed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
