package model

import (
	"time"

	"github.com/canonflow/backend-starter/pkg/helpers"
	"gorm.io/gorm"
)

type Permission struct {
	ID          string     `json:"id,omitempty" gorm:"primaryKey"`
	ActionID    string     `json:"action_id,omitempty" gorm:"column:action_id"`
	Action      *Action    `json:"action,omitempty"`
	Resource    string     `json:"resource,omitempty" gorm:"column:resource"`
	Description string     `json:"description,omitempty" gorm:"column:description"`
	CreatedAt   *time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	Roles       []Role     `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = helpers.GenerateUUIDV7()

	return nil
}
