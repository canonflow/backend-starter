package handler

import (
	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/core"
	"github.com/canonflow/backend-starter/internal/dto"
	usecase "github.com/canonflow/backend-starter/internal/usecase/user"
	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	LogPrefixSignUp    = "[USR-SU]"
	ErrCodeSignUpEmail = "[USR-SU-001]" + " BAD_REQUEST"

	LogPrefixSignIn             = "[USR-SI]"
	ErrCodeSignInInvalidPayload = "[USR-SI-001]" + " BAD_REQUEST"
	ErrCodeSignInEmail          = "[USR-SI-002]" + " NOT_FOUND"
	ErrCodeSignInPassword       = "[USR-SI-003]" + " NOT_FOUND"
	ErrCodeSignInAccessToken    = "[USR-SI-004]" + " INTERNAL_SERVER_ERROR"
)

type UserHandler struct {
	UserUsecase usecase.IUserUsecase
}

func NewUserHandler(userUsecase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) SignUp(ctx fiber.Ctx) error {
	var userDto dto.CreateUserRequest
	logger := core.LoggerFromContext(ctx.Context())

	if err := ctx.Bind().Body(&userDto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Info(LogPrefixSignUp+" Request payload parsed successfully",
		zap.String("email", userDto.Email),
	)

	// Check the username
	_, err := h.UserUsecase.FindBy(ctx.Context(), "email", userDto.Email, false)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(response.Error(ErrCodeSignUpEmail, "Email is already taken", "email"))
	}

	// Create user
	user, err := h.UserUsecase.Create(ctx.Context(), userDto.Email, userDto.Password)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(response.Success(
			user,
			nil,
		))
}

func (h *UserHandler) SignIn(ctx fiber.Ctx) error {
	var userDto dto.SignInRequest

	if err := ctx.Bind().Body(&userDto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, ErrCodeSignInInvalidPayload)
	}

	// Find by Email
	user, err := h.UserUsecase.FindBy(ctx.Context(), "email", userDto.Email, false)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).
			JSON(response.Error(ErrCodeSignInPassword, "Not Found", "credentials"))
	}

	// Check Password
	if !h.UserUsecase.Authenticate(user, userDto.Password) {
		return ctx.Status(fiber.StatusNotFound).
			JSON(response.Error(ErrCodeSignInPassword, "Not Found", "credentials"))
	}

	// Create Access Token
	tokenString, err := h.UserUsecase.CreateAccessToken(
		ctx.Context(),
		int(user.ID),
		user.Email,
		config.Get[string](config.AppKey),
		config.Get[uint](config.JWTDurationInMinute),
	)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(response.Error(ErrCodeSignInAccessToken, "Internal Server Error", "-"))
	}

	// Set Access Token to HTTP-Only Cookie
	wrapper := core.NewFiberContextWrapper(ctx)
	wrapper.SetToken(tokenString, config.Get[uint](config.JWTDurationInMinute))

	return ctx.Status(fiber.StatusOK).
		JSON(response.Success(
			user,
			nil,
		))
}
