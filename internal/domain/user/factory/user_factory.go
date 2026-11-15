package factory

import (
	"github.com/canonflow/backend-starter/internal/domain/user/repository"
	userUseCase "github.com/canonflow/backend-starter/internal/domain/user/usecase"
	"gorm.io/gorm"
)

func ResolveUserRepository(db *gorm.DB) repository.IUserRepository {
	return repository.NewUserRepository_MySQL(db)
}

func ResolveUserUseCaseV1(
	db *gorm.DB,
	userRepo repository.IUserRepository,
) userUseCase.IUserUsecase {
	return userUseCase.NewUserUsecaseV1(db, userRepo)
}
