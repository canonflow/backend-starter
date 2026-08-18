package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB    *gorm.DB
	App   *fiber.App
	Redis *redis.Client
}

func Bootstrap(cfg *BootstrapConfig) {
	// Setup Repo

	// Setup Usecase

	// Setup Handler

	// Setup Middleware

	// Setup Route

	// Init
}
