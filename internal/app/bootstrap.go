package app

import (
	"github.com/canonflow/backend-starter/internal/app/cache"
	"github.com/canonflow/backend-starter/internal/app/permission"
	accessFactory "github.com/canonflow/backend-starter/internal/domain/access/factory"
	userFactory "github.com/canonflow/backend-starter/internal/domain/user/factory"
	"github.com/canonflow/backend-starter/internal/http/delivery"
	"github.com/canonflow/backend-starter/internal/http/handler"
	"github.com/canonflow/backend-starter/internal/http/middleware"

	"github.com/gofiber/fiber/v3"
	fiberredis "github.com/gofiber/storage/redis/v3"
	redigo "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB    *gorm.DB
	App   *fiber.App
	Redis *redigo.Client
}

func Bootstrap(cfg *BootstrapConfig) {
	cacheRepo := cache.ResolveCacheRepositoryRedis()
	permissionAccessStorage := permission.ResolvePermissionRepository(cfg.DB, cfg.Redis)

	// Setup Repo
	userRepository := userFactory.ResolveUserRepository(cfg.DB)
	accessRepo := accessFactory.ResolveAccessRepository(cfg.DB)

	// Setup Usecase
	userUsecase := userFactory.ResolveUserUseCaseV1(cfg.DB, userRepository)
	accessUsecase := accessFactory.ResolveAccessUseCaseV1(cfg.DB, cacheRepo, accessRepo)

	// Setup Handler
	userHandler := handler.NewUserHandler(userUsecase, permissionAccessStorage)
	accessHandler := handler.NewRoleAccessHandler(accessUsecase)

	// Setup Middleware
	requestLoggerMiddleware := middleware.RequestLoggerMiddleware()
	api := cfg.App.Group("/api", requestLoggerMiddleware)

	// Setup Route
	throttleStorage := fiberredis.NewFromConnection(cfg.Redis)
	routeGroup := delivery.RouteGroup{}

	userRoute := delivery.NewUserRoute(api, userHandler, throttleStorage, permissionAccessStorage)
	accessRoute := delivery.NewRoleAccessRoute(api, accessHandler, throttleStorage, permissionAccessStorage)
	routeGroup.Register(userRoute, accessRoute)

	// Init
	routeGroup.Wire()
}
