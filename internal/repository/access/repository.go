package model

import (
	"context"

	"github.com/canonflow/backend-starter/internal/model"
	"github.com/canonflow/backend-starter/pkg/response"
	"gorm.io/gorm"
)

type IRoleAccessRepository interface {
	ListRole(context context.Context, pagination *response.Pagination, withPermission bool) ([]model.Role, error)
	FindRoleBy(context context.Context, column string, value any, withPermission bool) (*model.Role, error)
	CreateRole(db *gorm.DB, role *model.Role) error
	UpdateRole(db *gorm.DB, role *model.Role) error
	AddPermissionToRole(db *gorm.DB, role *model.Role, permissionId string) error

	ListPermission(context context.Context, pagination *response.Pagination, withAction bool) ([]model.Permission, error)
	FindPermissionBy(context context.Context, column string, value any, withAction bool) (*model.Permission, error)
	CreatePermission(db *gorm.DB, permission *model.Permission) error
	DeletePermission(db *gorm.DB, permission *model.Permission) error

	ListAction(context context.Context, pagination *response.Pagination, withPermission bool) ([]model.Action, error)
	FindActionBy(context context.Context, column string, value any, withPermission bool) (*model.Action, error)
	CreateAction(db *gorm.DB, action *model.Action) error
}
