package delivery

import (
	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/http/handler"
	"github.com/canonflow/backend-starter/internal/http/middleware"
	"github.com/gofiber/fiber/v3"

	fiberRedis "github.com/gofiber/storage/redis/v3"
)

type UserRoute struct {
	App             fiber.Router
	UserHandler     *handler.UserHandler
	ThrottleStorage *fiberRedis.Storage
}

func NewUserRoute(app fiber.Router, userHandler *handler.UserHandler, redisClient *fiberRedis.Storage) *UserRoute {
	return &UserRoute{
		App:             app,
		UserHandler:     userHandler,
		ThrottleStorage: redisClient,
	}
}

func (r *UserRoute) Setup() {
	userV1Path := r.App.Group("/v1/user")

	authMiddleware := middleware.AuthMiddleware(config.Get[string](config.AppKey))

	userV1Path.Get(
		"/me",
		middleware.Throttle(middleware.ThrottleConfig{
			MaxAttempts:  5,
			DecayMinutes: 1,
			KeyPrefix:    "me",
			Storage:      r.ThrottleStorage,
		}),
		authMiddleware,
		func(c fiber.Ctx) error {
			return c.SendString("Test Throttle")
		},
	)
	userV1Path.Post(
		"/signup",
		middleware.Throttle(middleware.ThrottleConfig{
			MaxAttempts:  5,
			DecayMinutes: 1,
			KeyPrefix:    "signup",
			Storage:      r.ThrottleStorage,
		}),
		r.UserHandler.SignUp,
	)
	userV1Path.Post(
		"/signin",
		middleware.Throttle(middleware.ThrottleConfig{
			MaxAttempts:  5,
			DecayMinutes: 1,
			KeyPrefix:    "signin",
			Storage:      r.ThrottleStorage,
		}),
		r.UserHandler.SignIn,
	)
}
