package model

import (
	internalRepo "github.com/canonflow/backend-starter/internal/repository"
	"gorm.io/gorm"
)

func NewRoleRepositoryFactory(db *gorm.DB, driver internalRepo.RepositoryDriver) IRoleAccessRepository {
	switch driver {
	case internalRepo.MySQLDriver:
		return newRoleRepository_MySQL(db)
	default:
		return newRoleRepository_MySQL(db)
	}
}
