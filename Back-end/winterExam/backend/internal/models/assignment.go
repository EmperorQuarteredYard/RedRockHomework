package models

import (
	"gorm.io/gorm"
	"time"
)

type Assignment struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Department  string         `json:"department" gorm:"type:enum('backend','frontend','sre','product','design','android','ios');not null"`
	CreatorID   uint64         `json:"creator_id" gorm:"not null"` // 发布者ID
	Deadline    time.Time      `json:"deadline"`
	AllowLate   bool           `json:"allow_late" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (a *Assignment) DeadlineString() string {
	return a.Deadline.Format("2006-01-02 15:04:05")
}
