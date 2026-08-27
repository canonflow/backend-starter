package model

import (
	"context"

	"github.com/canonflow/backend-starter/internal/app/scope"
	"github.com/canonflow/backend-starter/internal/model"
	"github.com/canonflow/backend-starter/pkg/response"
	"gorm.io/gorm"
)

type RoleAccessRepository_MySQL struct {
	DB *gorm.DB
}

// func newRoleRepository_MySQL(db *gorm.DB) IRoleRepository {
// 	return &RoleRepository_MySQL{
// 		DB: db,
// 	}
// }

func (r *RoleAccessRepository_MySQL) ListRole(context context.Context, pagination *response.Pagination, withPermission bool) ([]model.Role, error) {
	var roles []model.Role

	query := r.DB.WithContext(context).
		Scopes(scope.Paginate(roles, pagination, r.DB))

	if withPermission {
		query = query.Scopes(scope.RoleWithPermission)
	}

	err := query.Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *RoleAccessRepository_MySQL) FindRoleBy(context context.Context, column string, value any, withPermission bool) (*model.Role, error) {
	var role model.Role

	query := r.DB.WithContext(context)

	if withPermission {
		query = query.Scopes(scope.RoleWithPermission)
	}

	err := query.Where(map[string]any{
		column: value,
	}).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleAccessRepository_MySQL) CreateRole(db *gorm.DB, role *model.Role) error {
	return db.Create(role).Error
}

func (r *RoleAccessRepository_MySQL) UpdateRole(db *gorm.DB, role *model.Role) error {
	return db.Save(role).Error
}

func (r *RoleAccessRepository_MySQL) AddPermissionToRole(db *gorm.DB, role *model.Role, permissionId string) error {
	return db.Model(role).
		Association("Permissions").
		Append(&model.Permission{ID: permissionId})
}

func (r *RoleAccessRepository_MySQL) ListPermission(context context.Context, pagination *response.Pagination, withAction bool) ([]model.Permission, error) {
	var permissions []model.Permission

	query := r.DB.WithContext(context)

	if withAction {
		query = query.Scopes(scope.PermissionWithAction)
	}

	err := query.Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *RoleAccessRepository_MySQL) FindPermissionBy(context context.Context, column string, value any, withAction bool) (*model.Permission, error) {
	var permission model.Permission

	query := r.DB.WithContext(context)

	if withAction {
		query = query.Scopes(scope.PermissionWithAction)
	}

	err := query.Where(map[string]any{
		column: value,
	}).First(&permission).Error
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *RoleAccessRepository_MySQL) CreatePermission(db *gorm.DB, permission *model.Permission) error {
	return db.Create(permission).Error
}

func (r *RoleAccessRepository_MySQL) DeletePermission(db *gorm.DB, permission *model.Permission) error {
	return db.Save(permission).Error
}

func (r *RoleAccessRepository_MySQL) ListAction(context context.Context, pagination *response.Pagination, withPermission bool) ([]model.Action, error) {
	var actions []model.Action

	query := r.DB.WithContext(context)

	if withPermission {
		query = query.Scopes(scope.ActionWithPermission)
	}

	err := query.Find(&actions).Error
	if err != nil {
		return nil, err
	}

	return actions, nil
}

func (r *RoleAccessRepository_MySQL) FindActionBy(context context.Context, column string, value any, withPermission bool) (*model.Action, error) {
	var action model.Action

	query := r.DB.WithContext(context)

	if withPermission {
		query = query.Scopes(scope.ActionWithPermission)
	}

	err := query.Where(map[string]any{
		column: value,
	}).First(&action).Error
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func (r *RoleAccessRepository_MySQL) CreateAction(db *gorm.DB, action *model.Action) error {
	return db.Create(action).Error
}
