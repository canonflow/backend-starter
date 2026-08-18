package usecase

import (
	"context"
	"errors"

	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/model"
	userRepository "github.com/canonflow/backend-starter/internal/repository/user"
	"github.com/canonflow/backend-starter/pkg"
	"github.com/canonflow/backend-starter/pkg/response"
	"gorm.io/gorm"
)

type UserUsecaseV1 struct {
	DB             *gorm.DB
	UserRepository userRepository.IUserRepository
}

func NewUserUsecaseV1(db *gorm.DB, userRepository userRepository.IUserRepository) IUserUsecase {
	return &UserUsecaseV1{
		DB:             db,
		UserRepository: userRepository,
	}
}

func (u *UserUsecaseV1) CreateAccessToken(ctx context.Context, id int, email, key string, durationInMinutes uint) (string, error) {
	accessToken, err := pkg.CreateAccessToken(
		id,
		email,
		config.Get[string](config.AppKey),
		config.Get[uint](config.JWTDurationInMinute),
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (u *UserUsecaseV1) Authenticate(user *model.User, password string) bool {
	return pkg.CheckHash(password, user.Password)
}

func (u *UserUsecaseV1) List(ctx context.Context, limit, page int, sortBy, sort string, withTrashed bool) ([]model.User, *response.Pagination, error) {
	pagination := &response.Pagination{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		SortBy: sortBy,
	}

	users, err := u.UserRepository.List(
		ctx,
		pagination,
		withTrashed,
	)
	if err != nil {
		return nil, nil, err
	}

	return users, pagination, nil
}

func (u *UserUsecaseV1) Create(ctx context.Context, email, password string) (model.User, error) {
	// Find the user
	_, err := u.UserRepository.FindBy(ctx, "email", email, false)
	if err == nil {
		return model.User{}, errors.New("Email is already taken")
	}

	// Generate Hashed-password
	hashedPassword, err := pkg.Hash(password)
	if err != nil {
		return model.User{}, err
	}

	// Transaction
	tx := u.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return model.User{}, tx.Error
	}
	defer tx.Rollback()

	// Create user
	user := model.User{
		Email:     email,
		Password:  hashedPassword,
		DeletedAt: nil,
	}

	err = u.UserRepository.Create(ctx, tx, &user)
	if err != nil {
		return model.User{}, nil
	}

	if err := tx.Commit().Error; err != nil {
		return model.User{}, nil
	}

	return user, nil
}

func (u *UserUsecaseV1) Update(ctx context.Context, user *model.User) error {
	// Find the user
	_, err := u.UserRepository.FindBy(ctx, "email", user.Email, false)
	if err != nil {
		return err
	}

	// Update user
	tx := u.DB.
		WithContext(ctx).
		Begin()
	defer tx.Rollback()

	if err := u.UserRepository.Update(ctx, tx, user); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (u *UserUsecaseV1) Delete(ctx context.Context, user *model.User) error {
	// Find the user
	_, err := u.UserRepository.FindBy(ctx, "email", user.Email, false)
	if err != nil {
		return err
	}

	// Delete user
	tx := u.DB.
		WithContext(ctx).
		Begin()
	defer tx.Rollback()

	if err := u.UserRepository.Delete(ctx, tx, user); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
