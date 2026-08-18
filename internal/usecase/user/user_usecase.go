package usecase

import (
	"context"

	"github.com/canonflow/backend-starter/internal/model"
	"github.com/canonflow/backend-starter/pkg/response"
)

type IUserUsecase interface {
	CreateAccessToken(ctx context.Context, id int, email, key string, durationInMinutes uint) (string, error)
	Authenticate(user *model.User, password string) bool
	List(ctx context.Context, limit, page int, sortBy, sort string, withTrashed bool) ([]model.User, *response.Pagination, error)
	Create(ctx context.Context, email, password string) (model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
}
