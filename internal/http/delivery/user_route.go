package delivery

import (
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

	userV1Path.Get(
		"/",
		middleware.Throttle(middleware.ThrottleConfig{
			MaxAttempts:  5,
			DecayMinutes: 1,
			KeyPrefix:    "test",
			Storage:      r.ThrottleStorage,
		}),
		func(c fiber.Ctx) error {
			return c.SendString("Test Throttle")
		},
	)
	userV1Path.Post("/signup", r.UserHandler.SignUp)
	userV1Path.Post("/signin", r.UserHandler.SignIn)
}
