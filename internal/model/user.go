package model

import (
	"time"
)

type User struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Email     string     `gorm:"column:email;unique;not null" json:"email,omitempty"`
	Password  string     `gorm:"column:password;not null" json:"-"`
	DeletedAt *time.Time `gorm:"column:deleted_at;default:null" json:"deleted_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	Roles     []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}
