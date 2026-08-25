package model

import (
	"context"

	"github.com/canonflow/backend-starter/internal/app/scope"
	"github.com/canonflow/backend-starter/internal/model"
	"github.com/canonflow/backend-starter/pkg/response"
	"gorm.io/gorm"
)

type RoleRepository_MySQL struct {
	DB *gorm.DB
}

// func newRoleRepository_MySQL(db *gorm.DB) IRoleRepository {
// 	return &RoleRepository_MySQL{
// 		DB: db,
// 	}
// }

func (r *RoleRepository_MySQL) List(context context.Context, pagination *response.Pagination, withPermission bool) ([]model.Role, error) {
	var roles []model.Role

	query := r.DB.WithContext(context).
		Scopes(scope.Paginate(roles, pagination, r.DB))

	if withPermission {
		query.Preload("Permissions").
			Preload("Permissions.Actions", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name")
			})
	}

	err := query.Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}
func (r *RoleRepository_MySQL) FindBy(context context.Context, column string, value any, withTrashed bool) (*model.Role, error)
func (r *RoleRepository_MySQL) Create(db *gorm.DB, role *model.Role) error
func (r *RoleRepository_MySQL) Update(db *gorm.DB, role *model.Role) error
func (r *RoleRepository_MySQL) Delete(db *gorm.DB, role *model.Role) error
