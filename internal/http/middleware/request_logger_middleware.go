package middleware

import (
	"github.com/canonflow/backend-starter/internal/core"
	"github.com/gofiber/fiber/v3"
)

func RequestLoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := core.WithLogger(c.Context())
		c.SetContext(ctx) // Context with Zap Logger

		return c.Next()
	}
}
