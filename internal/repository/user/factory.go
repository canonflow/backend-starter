package repository

import (
	internalRepo "github.com/canonflow/backend-starter/internal/repository"
	"gorm.io/gorm"
)

func NewUserRepositoryFactory(db *gorm.DB, driver internalRepo.RepositoryDriver) IUserRepository {
	switch driver {
	case internalRepo.MySQLDriver:
		return newUserRepository_MySQL(db)
	default:
		return newUserRepository_MySQL(db)
	}
}
