package delivery

import (
	permissionRepository "github.com/canonflow/backend-starter/internal/app/permission/repository"
	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/http/handler"
	"github.com/canonflow/backend-starter/internal/http/middleware"
	"github.com/gofiber/fiber/v3"
	fiberRedis "github.com/gofiber/storage/redis/v3"
)

type RoleAccessRoute struct {
	App               fiber.Router
	RoleAccessHandler *handler.RoleAccessHandler
	ThrottleStorage   *fiberRedis.Storage
	PermissionAccess  permissionRepository.IPermissionAccess
}

func NewRoleAccessRoute(
	app fiber.Router,
	roleAccessHandler *handler.RoleAccessHandler,
	redisClient *fiberRedis.Storage,
	permissionAccess permissionRepository.IPermissionAccess,
) *RoleAccessRoute {
	return &RoleAccessRoute{
		App:               app,
		RoleAccessHandler: roleAccessHandler,
		ThrottleStorage:   redisClient,
		PermissionAccess:  permissionAccess,
	}
}

func (r *RoleAccessRoute) Setup() {
	authMiddleware := middleware.AuthMiddleware(config.Get[string](config.AppKey))
	roleV1Path := r.App.Group("/v1/role", authMiddleware)

	roleV1Path.Get(
		"/",
		// middleware.Throttle(middleware.ThrottleConfig{
		// 	MaxAttempts:  5,
		// 	DecayMinutes: 1,
		// 	KeyPrefix:    "role_route",
		// 	Storage:      r.ThrottleStorage,
		// }),
		r.RoleAccessHandler.ListOfRole,
	)
}
