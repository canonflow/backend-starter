package delivery

import (
	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/contract"
	"github.com/canonflow/backend-starter/internal/http/handler"
	"github.com/canonflow/backend-starter/internal/http/middleware"
	"github.com/gofiber/fiber/v3"

	fiberRedis "github.com/gofiber/storage/redis/v3"
)

type UserRoute struct {
	App              fiber.Router
	UserHandler      *handler.UserHandler
	ThrottleStorage  *fiberRedis.Storage
	PermissionAccess contract.IPermissionAccess
}

func NewUserRoute(
	app fiber.Router,
	userHandler *handler.UserHandler,
	redisClient *fiberRedis.Storage,
	permissionAccess contract.IPermissionAccess,
) *UserRoute {
	return &UserRoute{
		App:              app,
		UserHandler:      userHandler,
		ThrottleStorage:  redisClient,
		PermissionAccess: permissionAccess,
	}
}

func (r *UserRoute) Setup() {
	userV1Path := r.App.Group("/v1/user")

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

	authMiddleware := middleware.AuthMiddleware(config.Get[string](config.AppKey))
	authPath := userV1Path.Group("", authMiddleware)

	authPath.Get(
		"/me",
		middleware.Throttle(middleware.ThrottleConfig{
			MaxAttempts:  5,
			DecayMinutes: 1,
			KeyPrefix:    "me",
			Storage:      r.ThrottleStorage,
		}),
		r.UserHandler.Me,
	)

	authPath.Get(
		"/permissions",
		middleware.Throttle(middleware.ThrottleConfig{
			MaxAttempts:  10,
			DecayMinutes: 1,
			KeyPrefix:    "permissions",
			Storage:      r.ThrottleStorage,
		}),
		middleware.Permission(
			middleware.PermissionConfig{
				PermissionAccess: r.PermissionAccess,
			},
			middleware.OnResource("permission", "view"),
		),
		r.UserHandler.GetPermission,
	)

	authPath.Post(
		"signout",
		middleware.Throttle(middleware.ThrottleConfig{
			MaxAttempts:  2,
			DecayMinutes: 1,
			KeyPrefix:    "signout",
			Storage:      r.ThrottleStorage,
		}),
		r.UserHandler.SignOut,
	)
}
