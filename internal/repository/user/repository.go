package repository

import (
	"context"

	"github.com/canonflow/backend-starter/internal/model"
	"github.com/canonflow/backend-starter/pkg/response"
	"gorm.io/gorm"
)

type IUserRepository interface {
	List(context context.Context, pagination *response.Pagination, withTrashed bool) ([]model.User, error)
	FindBy(context context.Context, column string, value any, withTrashed bool) (*model.User, error)
	Create(context context.Context, db *gorm.DB, user *model.User) error
	Update(context context.Context, db *gorm.DB, user *model.User) error
	Delete(context context.Context, db *gorm.DB, user *model.User) error
}
