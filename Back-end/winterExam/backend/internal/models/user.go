package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string         `json:"username" gorm:"column:name;unique;not null"` // 对应原 Name 字段
	PasswordHash string         `json:"-" gorm:"column:password_hash;not null"`
	Nickname     string         `json:"nickname" gorm:"not null"`
	Role         string         `json:"role" gorm:"default:student"`
	Department   string         `json:"department" gorm:"type:enum('backend','frontend','sre','product','design','android','ios');not null"`
	Email        string         `json:"email"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) DepartmentLabel() string {
	dept, _ := GetDepartment(u.Department)
	return dept.DepartmentLabel
}
