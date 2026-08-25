package model

import (
	"time"

	"github.com/canonflow/backend-starter/pkg/helpers"
	"gorm.io/gorm"
)

type Action struct {
	ID        string     `json:"id,omitempty" gorm:"primarykey"`
	Name      string     `json:"name,omitempty" gorm:"name"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
}

func (a *Action) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = helpers.GenerateUUIDV7()

	return nil
}
