package app

import (
	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/internal/core/permission"
	"github.com/canonflow/backend-starter/internal/http/delivery"
	"github.com/canonflow/backend-starter/internal/http/handler"
	"github.com/canonflow/backend-starter/internal/http/middleware"
	internalRepo "github.com/canonflow/backend-starter/internal/repository"
	userRepo "github.com/canonflow/backend-starter/internal/repository/user"
	usecase "github.com/canonflow/backend-starter/internal/usecase/user"

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
	permissionAccessStorage := permission.NewPermissionAccessStorage(cfg.DB, cfg.Redis)

	// Setup Repo
	userRepository := userRepo.NewUserRepositoryFactory(cfg.DB, internalRepo.GetDriver(config.GetOrDefault(config.DBDriver, "mysql")))

	// Setup Usecase
	userUsecase := usecase.NewUserUsecaseV1(cfg.DB, userRepository)

	// Setup Handler
	userHandler := handler.NewUserHandler(userUsecase, permissionAccessStorage)

	// Setup Middleware
	requestLoggerMiddleware := middleware.RequestLoggerMiddleware()
	api := cfg.App.Group("/api", requestLoggerMiddleware)

	// Setup Route
	throttleStorage := fiberredis.NewFromConnection(cfg.Redis)
	routeGroup := delivery.RouteGroup{}
	userRoute := delivery.NewUserRoute(api, userHandler, throttleStorage, permissionAccessStorage)
	routeGroup.Register(userRoute)

	// Init
	routeGroup.Wire()
}
