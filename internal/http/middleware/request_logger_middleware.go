package middleware

import (
	"time"

	"github.com/canonflow/backend-starter/internal/core"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := core.WithLogger(c.Context())
		c.SetContext(ctx) // Context with Zap Logger

		start := time.Now()
		err := c.Next()

		duration := time.Since(start)
		logger := core.LoggerFromContext(c.Context())
		logger.Info("[REQUEST LOGGER] request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Int64("duration_ms", duration.Microseconds()),
			zap.String("method", c.IP()),
		)

		return err
	}
}
