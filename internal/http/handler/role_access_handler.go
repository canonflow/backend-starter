package handler

import (
	"fmt"

	appContext "github.com/canonflow/backend-starter/internal/app/context"
	"github.com/canonflow/backend-starter/internal/domain/access/dto"
	usecase "github.com/canonflow/backend-starter/internal/domain/access/usecase"
	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	LogPrefixListRole = "[RA-LR]"
	ErrCodeList       = "[RA-LR-001] " + "INTERNAL_SERVER_ERROR"
)

type RoleAccessHandler struct {
	RoleAccessUsecase usecase.IRoleAccessUsecase
}

func NewRoleAccessHandler(roleAccessHandler usecase.IRoleAccessUsecase) *RoleAccessHandler {
	return &RoleAccessHandler{
		RoleAccessUsecase: roleAccessHandler,
	}
}

func (h *RoleAccessHandler) ListOfRole(ctx fiber.Ctx) error {
	payload := new(dto.ListRoleRequest)

	if err := ctx.Bind().Query(payload); err != nil {
		appContext.LoggerFromContext(ctx.Context()).
			Info(LogPrefixListRole+" Error parsing query", zap.String("error", err.Error()))
		return err
	}

	appContext.LoggerFromContext(ctx.Context()).
		Info(LogPrefixListRole+" Payload",
			zap.Any("payload", payload),
			zap.String("request_uri", string(ctx.Request().URI().QueryString())),
		)
	roles, pagination, err := h.RoleAccessUsecase.ListRole(
		ctx.Context(),
		payload.Limit,
		payload.Page,
		payload.SearchColumn,
		payload.SearchValue,
		payload.SortBy,
		payload.Sort,
		true,
	)
	if err != nil {
		appContext.LoggerFromContext(ctx.Context()).
			Info(
				fmt.Sprintf(LogPrefixListRole+" Error while getting role from usecase"),
				zap.String("error", err.Error()),
			)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(response.Error(ErrCodeList, "Internal Server Error", "-"))
	}

	r := response.Success(roles, &response.Meta{
		Pagination: pagination,
	})

	return ctx.Status(fiber.StatusOK).
		JSON(r)
}
