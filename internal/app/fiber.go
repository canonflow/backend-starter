package app

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/canonflow/backend-starter/internal/config"
	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		// Prefork ListenConfig when serve the fiber
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			codeMessage := "[0000] INTERNAL_SERVER_ERROR"
			message := "Internal Server Error"
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
				codeMessage = message
			}

			return ctx.Status(code).
				JSON(response.Error(codeMessage, message, "_"))
		},
	})

	// Panic Recovery
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx fiber.Ctx, e interface{}) {
			logger := LoggerFromContext(ctx.Context())
			logger.Error(fmt.Sprintf("[0000] Panic Occured: %v\n", e))
		},
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetOrDefault(config.CORSAllowOrigins, []string{"*"}),
		AllowHeaders:     config.GetOrDefault(config.CORSAllowHeaders, []string{}),
		AllowMethods:     config.GetOrDefault(config.CORSAllowMethods, []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "QUERY"}),
		AllowCredentials: config.GetOrDefault(config.CORSAllowCredentials, true),
	}))

	return app
}
