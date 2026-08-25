package handler

import (
	"errors"
	"fmt"

	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/contract"
	"github.com/canonflow/backend-starter/internal/core"
	"github.com/canonflow/backend-starter/internal/dto"
	usecase "github.com/canonflow/backend-starter/internal/usecase/user"
	"github.com/canonflow/backend-starter/pkg"
	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	LogPrefixSignUp    = "[USR-SU]"
	ErrCodeSignUpEmail = "[USR-SU-001]" + " BAD_REQUEST"

	LogPrefixSignIn             = "[USR-SI]"
	ErrCodeSignInInvalidPayload = "[USR-SI-001]" + " BAD_REQUEST"
	ErrCodeSignInEmail          = "[USR-SI-002]" + " NOT_FOUND"
	ErrCodeSignInPassword       = "[USR-SI-003]" + " NOT_FOUND"
	ErrCodeSignInAccessToken    = "[USR-SI-004]" + " INTERNAL_SERVER_ERROR"

	LogPrefixMe = "[USR-ME]"

	LogPrefixGetPermission = "[USR-GP]"
	ErrCodeGetPermission   = "[USR-GP-01]" + "INTERNAL_SERVER_ERROR"
)

type UserHandler struct {
	UserUsecase      usecase.IUserUsecase
	PermissionAccess contract.IPermissionAccess
}

func NewUserHandler(userUsecase usecase.IUserUsecase, permissionAccess contract.IPermissionAccess) *UserHandler {
	return &UserHandler{
		UserUsecase:      userUsecase,
		PermissionAccess: permissionAccess,
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

	// Check the email
	_, err := h.UserUsecase.FindBy(ctx.Context(), "email", userDto.Email, false)
	if err == nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(response.Error(ErrCodeSignUpEmail, "Email is already taken", "email"))
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(LogPrefixSignUp+" Failed to check existing email",
			zap.String("email", userDto.Email),
			zap.Error(err),
		)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(response.Error("INTERNAL_SERVER_ERROR", "Something went wrong", "-"))
	}

	logger.Info(LogPrefixSignUp+" Creating user...",
		zap.String("email", userDto.Email),
	)

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

func (h *UserHandler) SignOut(ctx fiber.Ctx) error {
	ctxWrapper := core.NewFiberContextWrapper(ctx)

	ctxWrapper.DeleteToken()

	return ctx.Status(fiber.StatusOK).
		JSON(response.Success(
			"success",
			nil,
		))
}

func (h *UserHandler) Me(ctx fiber.Ctx) error {
	ctxWrapper := core.NewFiberContextWrapper(ctx)
	logger := core.LoggerFromContext(ctx.Context())

	userId, ok := core.GetLocal[int](ctxWrapper, pkg.JWTUserIDKey)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(response.Error("UNAUTHORIZED_ACCESS", "Unauthorized", "-"))
	}

	user, err := h.UserUsecase.FindBy(ctx.Context(), "id", userId, false)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(response.Error("UNAUTHORIZED_ACCESS", "Unauthorized", "-"))
	}

	logger.Info(LogPrefixSignUp+" Parsing user information from token",
		zap.Int("id", userId),
	)

	return ctx.Status(fiber.StatusOK).
		JSON(response.Success(
			user,
			nil,
		))
}

func (h *UserHandler) GetPermission(ctx fiber.Ctx) error {
	ctxWrapper := core.NewFiberContextWrapper(ctx)
	logger := core.LoggerFromContext(ctx.Context())

	userId, ok := core.GetLocal[int](ctxWrapper, pkg.JWTUserIDKey)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(response.Error("UNAUTHORIZED_ACCESS", "Unauthorized", "-"))
	}

	permissions, err := h.PermissionAccess.GetUserPermissions(ctx.Context(), userId)
	if err != nil {
		logger.Info(LogPrefixGetPermission + fmt.Sprintf(" Error while getting permissions: %v\n", err.Error()))
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(response.Error(ErrCodeGetPermission, "Internal Server Error", "-"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.Success(
			permissions,
			nil,
		))
}
