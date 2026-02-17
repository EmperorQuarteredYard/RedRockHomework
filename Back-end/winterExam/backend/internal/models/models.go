package models

import "gorm.io/gorm"

const ()

type User struct {
	ID           uint64         `json:"id" gorm:"primaryKey;unique;not null;autoIncrement"`
	Nickname     string         `json:"nickname"`
	PasswordHash string         `json:"passwordHash"`
	Name         string         `json:"name"`
	Role         string         `json:"role"`
	Section      string         `json:"section"`
	Delete       gorm.DeletedAt `gorm:"index"`
}

type Submit struct {
	ID           uint64         `json:"id" gorm:"primaryKey;unique;not null;autoIncrement"`
	Score        int            `json:"score"`
	AssignmentID uint64         `json:"assignment_id"`
	Delete       gorm.DeletedAt `gorm:"index"`
}

type Assignment struct {
	ID     uint64         `json:"id" gorm:"primaryKey;unique;not null;autoIncrement"`
	Delete gorm.DeletedAt `gorm:"index"`
}
