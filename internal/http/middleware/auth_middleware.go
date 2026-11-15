package middleware

import (
	"errors"

	appContext "github.com/canonflow/backend-starter/internal/app/context"
	appVars "github.com/canonflow/backend-starter/internal/app/vars"
	"github.com/canonflow/backend-starter/pkg"
	"github.com/canonflow/backend-starter/pkg/response"
	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		ctxWrapper := appContext.NewFiberContextWrapper(ctx)

		token, err := ctxWrapper.GetToken()

		if err != nil || token == "" {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(response.Error("UNAUTHORIZED_ACCESS", "Unauthorized", "-"))
		}

		claims, err := pkg.ParseToken(token, jwtSecret)
		if err != nil {
			if errors.Is(err, pkg.ErrTokenExpired) {
				return ctx.Status(fiber.StatusUnauthorized).
					JSON(response.Error("TOKEN_EXPIRED", "Session expired, please log in again", "-"))
			}
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(response.Error("UNAUTHORIZED_ACCESS", "Unauthorized", "-"))
		}

		userID, err := pkg.GetUserIDFromClaims(claims)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(response.Error("UNAUTHORIZED_ACCESS", "Unauthorized", "-"))
		}

		email, err := pkg.GetEmailFromClaims(claims)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).
				JSON(response.Error("UNAUTHORIZED_ACCESS", "Unauthorized", "-"))
		}

		// Make user info available to downstream handlers
		appContext.SetLocal(ctxWrapper, appVars.JWTUserIDKey, userID)
		appContext.SetLocal(ctxWrapper, appVars.JWTEmailKey, email)

		return ctx.Next()
	}
}
