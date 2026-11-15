package middleware

import (
	"time"

	appContext "github.com/canonflow/backend-starter/internal/app/context"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := appContext.WithLogger(c.Context())
		c.SetContext(ctx) // Context with Zap Logger

		start := time.Now()
		err := c.Next()

		duration := time.Since(start)
		logger := appContext.LoggerFromContext(c.Context())
		logger.Info("[REQUEST LOGGER] request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Int64("duration_ms", duration.Microseconds()),
			zap.String("ip", c.IP()),
			zap.String("user-agent", c.UserAgent()),
		)

		return err
	}
}
