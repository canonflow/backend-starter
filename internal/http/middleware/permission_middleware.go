package middleware

import (
	"fmt"

	"github.com/canonflow/backend-starter/internal/contract"
	"github.com/canonflow/backend-starter/internal/core"
	"github.com/canonflow/backend-starter/pkg"
	"github.com/canonflow/backend-starter/pkg/helpers"
	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type on struct {
	Resource string
	Action   string
}

func (o *on) Get() string {
	return fmt.Sprintf("%s.%s", o.Resource, o.Action)
}

type PermissionConfig struct {
	PermissionAccess contract.IPermissionAccess
}

func OnResource(resource, action string) on {
	return on{
		Resource: resource,
		Action:   action,
	}
}

func Permission(cfg PermissionConfig, resource on, resources ...on) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		data := []string{resource.Get()}

		for _, r := range resources {
			data = append(data, r.Get())
		}

		data = helpers.RemoveDuplicates(data)

		ctxWrapper := core.NewFiberContextWrapper(ctx)
		userId, ok := core.GetLocal[int](ctxWrapper, pkg.JWTUserIDKey)

		if !ok {
			return ctx.Status(fiber.StatusForbidden).
				JSON(response.Error("FORBIDDEN", "Forbidden", "-"))
		}

		pass, err := cfg.PermissionAccess.HasPermission(
			ctx.Context(),
			userId,
			data,
		)

		if !pass || err != nil {
			return ctx.Status(fiber.StatusForbidden).
				JSON(response.Error("FORBIDDEN", "Forbidden", "-"))
		}

		return ctx.Next()
	}
}
