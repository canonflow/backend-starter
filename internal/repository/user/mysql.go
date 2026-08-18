package repository

import (
	"context"
	"time"

	"github.com/canonflow/backend-starter/internal/app/scope"
	"github.com/canonflow/backend-starter/internal/model"
	"github.com/canonflow/backend-starter/pkg/response"
	"gorm.io/gorm"
)

type UserRepository_MySQL struct {
	DB *gorm.DB
}

func newUserRepository_MySQL(db *gorm.DB) IUserRepository {
	return &UserRepository_MySQL{
		DB: db,
	}
}

func (r *UserRepository_MySQL) List(context context.Context, pagination *response.Pagination, withTrashed bool) ([]model.User, error) {
	var users []model.User

	query := r.DB.WithContext(context).
		Scopes(scope.Paginate(users, pagination, r.DB))

	if !withTrashed {
		query = query.Scopes(scope.WithoutTrashed)
	}

	err := query.
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository_MySQL) FindBy(context context.Context, column string, value any, withTrashed bool) (*model.User, error) {
	var user model.User

	query := r.DB.WithContext(context).
		Where(map[string]any{
			column: value,
		})

	if !withTrashed {
		query = query.Scopes(scope.WithoutTrashed)
	}

	err := query.
		Find(&user).
		Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository_MySQL) Create(context context.Context, db *gorm.DB, user *model.User) error {
	return db.WithContext(context).
		Create(user).Error
}

func (r *UserRepository_MySQL) Update(context context.Context, db *gorm.DB, user *model.User) error {
	return db.WithContext(context).
		Save(user).Error
}

func (r *UserRepository_MySQL) Delete(context context.Context, db *gorm.DB, user *model.User) error {
	now := time.Now()
	user.DeletedAt = &now
	return db.WithContext(context).
		Save(user).Error
}
