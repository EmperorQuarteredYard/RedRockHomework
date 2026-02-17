package models

import (
	"time"

	"gorm.io/gorm"
)

const ()

type User struct {
	ID         uint64 `json:"id" Gorm:"primaryKey;unique;not null;autoIncrement"`
	Department string `json:"department"`

	Nickname     string `json:"nickname"`
	PasswordHash string `json:"passwordHash"`
	Name         string `json:"name"`
	Role         string `json:"role"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Delete    gorm.DeletedAt `Gorm:"index"`
}

// Submit 提交的作业，之后应当完善存储提交的内容
type Submit struct {
	ID         uint64 `json:"id" Gorm:"primaryKey;unique;not null;autoIncrement"`
	Department string `json:"department"`

	Score        int    `json:"score"`
	AssignmentID uint64 `json:"assignment_id"`
	SubmitterID  uint64 `json:"submitter_id"`
	Comment      string `json:"comment"`
	Excellent    bool   `json:"excellent" Gorm:"default:false"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Delete    gorm.DeletedAt `Gorm:"index"`
}

type Assignment struct {
	ID         uint64 `json:"id" Gorm:"primaryKey;unique;not null;autoIncrement"`
	Department string `json:"department"`

	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Deadline               time.Time `json:"deadline"`
	LateSubmissionStrategy string    `json:"late_submission_strategy"`
	PublishBy              uint64    `json:"publish_by"`
	UpdateBy               uint64    `json:"update_by"`
	UpdaterName            string    `json:"updater_name"`
	PublisherName          string    `json:"publisher_name"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Delete    gorm.DeletedAt `Gorm:"index"`
}
