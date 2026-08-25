package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/canonflow/backend-starter/internal/contract"
	"github.com/canonflow/backend-starter/internal/core"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	PermissionCacheKeyPrefix = "user:permissions:"
	permissionCacheTTL       = 15 * time.Minute
)

type PermissionAccessStorage struct {
	DB  *gorm.DB
	Rdb *redis.Client
}

func permissionCacheKey(userID int) string {
	return fmt.Sprintf("%s%d", PermissionCacheKeyPrefix, userID)
}

func NewPermissionAccessStorage(db *gorm.DB, rdb *redis.Client) contract.IPermissionAccess {
	return &PermissionAccessStorage{
		DB:  db,
		Rdb: rdb,
	}
}

func (c *PermissionAccessStorage) HasPermission(ctx context.Context, userID int, resources []string) (bool, error) {
	permissions, err := c.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	permSet := make(map[string]struct{}, len(permissions))
	for _, p := range permissions {
		permSet[p] = struct{}{}
	}

	for _, r := range resources {
		if _, ok := permSet[r]; ok {
			return true, nil
		}
	}

	return false, nil
}

func (c *PermissionAccessStorage) GetUserPermissions(ctx context.Context, userID int) ([]string, error) {
	log := core.LoggerFromContext(ctx)

	key := permissionCacheKey(userID)

	cached, err := c.Rdb.Get(ctx, key).Result()
	if err == nil {
		var permissions []string
		if unmarshalErr := sonic.Unmarshal([]byte(cached), &permissions); unmarshalErr == nil {
			return permissions, nil
		}
	}

	var permissions []string
	err = c.DB.WithContext(ctx).Table("users").
		Select("DISTINCT CONCAT(permissions.resource, '.', actions.name)").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN actions ON actions.id = permissions.action_id").
		Where("users.id = ?", userID).
		Pluck("CONCAT(permissions.resource, '.', actions.name)", &permissions).Error
	if err != nil {
		return nil, err
	}

	// 3. Populate cache (best-effort — don't fail the request if Redis write fails)
	if data, marshalErr := sonic.Marshal(permissions); marshalErr == nil {
		if setErr := c.Rdb.Set(ctx, key, data, permissionCacheTTL).Err(); setErr != nil {
			log.Info(fmt.Sprintf("redis set error for %s: %v", key, setErr))
		}
	}

	return permissions, nil
}

func (c *PermissionAccessStorage) InvalidateUserPermissions(ctx context.Context, userID int) error {
	return c.Rdb.Del(ctx, permissionCacheKey(userID)).Err()
}
