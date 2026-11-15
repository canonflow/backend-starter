package cache

import (
	"github.com/canonflow/backend-starter/internal/app/cache/repository"
)

func ResolveCacheRepositoryRedis() repository.ICache {
	return repository.NewRedis()
}
