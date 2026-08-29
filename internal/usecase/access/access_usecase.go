package usecase

import (
	"context"

	"github.com/canonflow/backend-starter/internal/model"
	"github.com/canonflow/backend-starter/pkg/response"
)

type IRoleAccessUsecase interface {
	ListRole(context context.Context, limit, page int, sortBy, sort string, withPermission bool) ([]model.Role, *response.Pagination, error)
	FindRoleBy(context context.Context, column string, value any, withPermission bool) (*model.Role, error)
	CreateRole(context context.Context, name, description string) (model.Role, error)
	UpdateRole(context context.Context, role *model.Role) error
	AddPermissionToRole(context context.Context, role *model.Role, permissionId string) error

	ListPermission(context context.Context, limit, page int, sortBy, sort string, withAction bool) ([]model.Permission, *response.Pagination, error)
	FindPermissionBy(context context.Context, column string, value any, withAction bool) (*model.Permission, error)
	CreatePermission(context context.Context, actionId, resource, description string) (model.Permission, error)
	DeletePermission(context context.Context, permission *model.Permission) error

	ListAction(context context.Context, limit, page int, sortBy, sort string, withPermission bool) ([]model.Action, *response.Pagination, error)
	FindActionBy(context context.Context, column string, value any, withPermission bool) (*model.Action, error)
	CreateAction(context context.Context, name string) (model.Action, error)
}
