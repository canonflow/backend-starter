package factory

import (
	cacheRepo "github.com/canonflow/backend-starter/internal/app/cache/repository"
	"github.com/canonflow/backend-starter/internal/domain/access/repository"
	"github.com/canonflow/backend-starter/internal/domain/access/usecase"
	"gorm.io/gorm"
)

func ResolveAccessRepository(db *gorm.DB) repository.IRoleAccessRepository {
	return repository.NewwRoleRepository_MySQL(db)
}

func ResolveAccessUseCaseV1(
	db *gorm.DB,
	cacheRepo cacheRepo.ICache,
	accessRepo repository.IRoleAccessRepository,
) usecase.IRoleAccessUsecase {
	return usecase.NewRoleAccessUsecaseV1(db, cacheRepo, accessRepo)
}
