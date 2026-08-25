package model

import (
	"time"
)

type Role struct {
	ID          int           `json:"id,omitempty" gorm:"primaryKey"`
	Name        string        `json:"name,omitempty" gorm:"column:name;not null"`
	Description string        `json:"description,omitempty" gorm:"column:description"`
	CreatedAt   *time.Time    `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time    `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	Users       []User        `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}
