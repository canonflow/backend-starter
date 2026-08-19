package delivery

import (
	"github.com/canonflow/backend-starter/internal/http/handler"
	"github.com/gofiber/fiber/v3"
)

type UserRoute struct {
	App         fiber.Router
	UserHandler *handler.UserHandler
}

func NewUserRoute(app fiber.Router, userHandler *handler.UserHandler) *UserRoute {
	return &UserRoute{
		App:         app,
		UserHandler: userHandler,
	}
}

func (r *UserRoute) Setup() {
	userPath := r.App.Group("/v1/users")

	userPath.Post("/signup", r.UserHandler.SignUp)
	userPath.Post("/signin", r.UserHandler.SignIn)
}
