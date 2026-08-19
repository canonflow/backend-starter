package middleware

import (
	"github.com/canonflow/backend-starter/internal/app"
	"github.com/gofiber/fiber/v3"
)

func RequestLoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := app.WithLogger(c.Context())
		c.SetContext(ctx) // Context with Zap Logger

		return c.Next()
	}
}
