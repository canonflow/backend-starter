package middleware

import (
	"fmt"
	"time"

	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type ThrottleConfig struct {
	MaxAttempts  int
	DecayMinutes float64
	KeyPrefix    string                   // namespace key per route
	KeyGenerator func(c fiber.Ctx) string // Default use IP
	Storage      fiber.Storage
}

func Throttle(cfg ThrottleConfig) fiber.Handler {
	keyGen := cfg.KeyGenerator

	if keyGen == nil {
		keyGen = func(c fiber.Ctx) string {
			return c.IP()
		}
	}

	prefix := cfg.KeyPrefix
	if prefix == "" {
		prefix = "default"
	}

	return limiter.New(limiter.Config{
		Max:        cfg.MaxAttempts,
		Expiration: time.Duration(cfg.DecayMinutes * float64(time.Minute)),
		Storage:    cfg.Storage,
		KeyGenerator: func(c fiber.Ctx) string {
			return fmt.Sprintf("throttle:%s:%s", prefix, keyGen(c))
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(response.Error("TOO_MANY_REQUESTS", "Too many attempts. Please try again later.", "-"))
		},
		LimiterMiddleware: limiter.SlidingWindow{}, // avoids the fixed-window burst problem
	})
}
