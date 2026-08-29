package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/bytedance/sonic"
	"github.com/canonflow/backend-starter/internal/model"
	accessRepo "github.com/canonflow/backend-starter/internal/repository/access"
	"github.com/canonflow/backend-starter/pkg/helpers"
	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	RoleCacheKeyPrefix = "role:list:"
	RoleCacheTTL       = 5 * time.Minute

	PermissionCacheKeyPrefix = "permission:list"
	PerissionCacheTTL        = 5 * time.Minute

	ActionCacheKeyPrefix = "action:list"
	ActionCacheTTL       = 5 * time.Minute

	ErrAddPermissionToRoleExist = "The role is already have the permission"
)

type RoleAccessUsecaseV1 struct {
	DB                   *gorm.DB
	Rdb                  *redis.Client
	RoleAccessRepository accessRepo.IRoleAccessRepository
}

type cachedRoleList struct {
	Data       []model.Role         `json:"data"`
	Pagination *response.Pagination `json:"pagination"`
}

type cachePermissionList struct {
	Data       []model.Permission   `json:"data"`
	Pagination *response.Pagination `json:"pagination"`
}

type cacheActionList struct {
	Data       []model.Action       `json:"data"`
	Pagination *response.Pagination `json:"pagination"`
}

func NewRoleAccessUsecaseV1(db *gorm.DB, rdb *redis.Client, roleAccessRepository accessRepo.IRoleAccessRepository) IRoleAccessUsecase {
	return &RoleAccessUsecaseV1{
		DB:                   db,
		Rdb:                  rdb,
		RoleAccessRepository: roleAccessRepository,
	}
}

func buildCacheListKey(prefix string, limit, page int, sortBy, sort string) string {
	s, _ := sonic.Marshal(map[string]any{
		"limit":   limit,
		"page":    page,
		"sort_by": sortBy,
		"sort":    sort,
	})

	return prefix + string(s)
}

func (u *RoleAccessUsecaseV1) ListRole(context context.Context, limit, page int, sortBy, sort string, withPermission bool) ([]model.Role, *response.Pagination, error) {
	key := buildCacheListKey(RoleCacheKeyPrefix, limit, page, sortBy, sort)

	cached, err := u.Rdb.Get(context, key).Result()
	if err == nil {
		var result cachedRoleList

		if unmarshalErr := sonic.Unmarshal([]byte(cached), &result); unmarshalErr == nil {
			return result.Data, result.Pagination, nil
		}
	}

	pagination := response.Pagination{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		SortBy: sortBy,
	}

	roles, err := u.RoleAccessRepository.ListRole(
		context,
		&pagination,
		withPermission,
	)
	if err != nil {
		return nil, nil, err
	}

	if b, err := sonic.Marshal(cachedRoleList{Data: roles, Pagination: &pagination}); err == nil {
		_ = u.Rdb.Set(context, key, b, RoleCacheTTL)
	}

	return roles, &pagination, nil
}

func (u *RoleAccessUsecaseV1) FindRoleBy(context context.Context, column string, value any, withPermission bool) (*model.Role, error) {
	role, err := u.RoleAccessRepository.FindRoleBy(context, column, value, withPermission)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (u *RoleAccessUsecaseV1) CreateRole(context context.Context, name, description string) (model.Role, error) {
	tx := u.DB.WithContext(context).Begin()

	if tx.Error != nil {
		return model.Role{}, tx.Error
	}
	defer tx.Rollback()

	role := model.Role{
		Name:        name,
		Description: description,
	}

	if err := u.RoleAccessRepository.CreateRole(tx, &role); err != nil {
		return model.Role{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.Role{}, err
	}

	return role, nil
}

func (u *RoleAccessUsecaseV1) UpdateRole(context context.Context, role *model.Role) error {
	_, err := u.RoleAccessRepository.FindRoleBy(context, "id", role.ID, false)
	if err != nil {
		return err
	}

	// Update
	tx := u.DB.
		WithContext(context).
		Begin()

	if err := u.RoleAccessRepository.UpdateRole(tx, role); err != nil {
		return err
	}

	return nil
}

func (u *RoleAccessUsecaseV1) AddPermissionToRole(context context.Context, role *model.Role, permissionId string) error {
	r, err := u.RoleAccessRepository.FindRoleBy(context, "id", role.ID, true)
	if err != nil {
		return err
	}

	permissions := make([]string, 0, len(r.Permissions))

	for _, p := range r.Permissions {
		permissions = append(permissions, p.ID)
	}

	if helpers.SliceContains(permissions, permissionId) {
		return errors.New(ErrAddPermissionToRoleExist)
	}

	tx := u.DB.
		WithContext(context).
		Begin()
	if err := u.RoleAccessRepository.AddPermissionToRole(tx, role, permissionId); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (u *RoleAccessUsecaseV1) ListPermission(context context.Context, limit, page int, sortBy, sort string, withAction bool) ([]model.Permission, *response.Pagination, error) {
	key := buildCacheListKey(PermissionCacheKeyPrefix, limit, page, sortBy, sort)

	cached, err := u.Rdb.Get(context, key).Result()
	if err == nil {
		var result cachePermissionList

		if unmarshalErr := sonic.Unmarshal([]byte(cached), &result); unmarshalErr == nil {
			return result.Data, result.Pagination, nil
		}
	}

	pagination := response.Pagination{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		SortBy: sortBy,
	}

	permissions, err := u.RoleAccessRepository.ListPermission(
		context,
		&pagination,
		withAction,
	)
	if err != nil {
		return nil, nil, err
	}

	if b, err := sonic.Marshal(cachePermissionList{Data: permissions, Pagination: &pagination}); err == nil {
		_ = u.Rdb.Set(context, key, b, PerissionCacheTTL)
	}

	return permissions, &pagination, nil
}

func (u *RoleAccessUsecaseV1) FindPermissionBy(context context.Context, column string, value any, withAction bool) (*model.Permission, error) {
	permission, err := u.RoleAccessRepository.FindPermissionBy(context, column, value, withAction)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (u *RoleAccessUsecaseV1) CreatePermission(context context.Context, actionId, resource, description string) (model.Permission, error) {
	tx := u.DB.WithContext(context).Begin()

	if tx.Error != nil {
		return model.Permission{}, tx.Error
	}
	defer tx.Rollback()

	permission := model.Permission{
		ActionID:    actionId,
		Resource:    resource,
		Description: description,
	}

	if err := u.RoleAccessRepository.CreatePermission(tx, &permission); err != nil {
		return model.Permission{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.Permission{}, err
	}

	return permission, nil
}

func (u *RoleAccessUsecaseV1) DeletePermission(context context.Context, permission *model.Permission) error {
	tx := u.DB.WithContext(context).Begin()

	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := u.RoleAccessRepository.DeletePermission(tx, permission); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (u *RoleAccessUsecaseV1) ListAction(context context.Context, limit, page int, sortBy, sort string, withPermission bool) ([]model.Action, *response.Pagination, error) {
	key := buildCacheListKey(ActionCacheKeyPrefix, limit, page, sortBy, sort)

	cached, err := u.Rdb.Get(context, key).Result()
	if err == nil {
		var result cacheActionList

		if unmarshalErr := sonic.Unmarshal([]byte(cached), &result); unmarshalErr == nil {
			return result.Data, result.Pagination, nil
		}
	}

	pagination := response.Pagination{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		SortBy: sortBy,
	}

	actions, err := u.RoleAccessRepository.ListAction(
		context,
		&pagination,
		withPermission,
	)
	if err != nil {
		return nil, nil, err
	}

	if b, err := sonic.Marshal(cacheActionList{Data: actions, Pagination: &pagination}); err == nil {
		_ = u.Rdb.Set(context, key, b, ActionCacheTTL)
	}

	return actions, &pagination, nil
}

func (u *RoleAccessUsecaseV1) FindActionBy(context context.Context, column string, value any, withPermission bool) (*model.Action, error) {
	action, err := u.RoleAccessRepository.FindActionBy(context, column, value, withPermission)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (u *RoleAccessUsecaseV1) CreateAction(context context.Context, name string) (model.Action, error) {
	tx := u.DB.WithContext(context).Begin()

	if tx.Error != nil {
		return model.Action{}, tx.Error
	}
	defer tx.Rollback()

	action := model.Action{
		Name: name,
	}

	if err := u.RoleAccessRepository.CreateAction(tx, &action); err != nil {
		return model.Action{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.Action{}, err
	}

	return action, nil
}
