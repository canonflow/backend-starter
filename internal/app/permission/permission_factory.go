package permission

import (
	"github.com/canonflow/backend-starter/internal/app/permission/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ResolvePermissionRepository(
	db *gorm.DB,
	rdb *redis.Client,
) repository.IPermissionAccess {
	return repository.NewPermissionAccessStorage(db, rdb)
}
