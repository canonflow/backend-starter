package model

import "time"

type User struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string     `gorm:"column:email;unique;not null" json:"email"`
	Password  string     `gorm:"column:password;not null" json:"-"`
	DeletedAt *time.Time `gorm:"column:deleted_at;default:null" json:"deleted_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`
}
