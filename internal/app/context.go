package app

import (
	"context"
	"time"

	"github.com/canonflow/backend-starter/pkg"
	"github.com/canonflow/backend-starter/pkg/helpers"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type loggerKey struct{}

func WithLogger(ctx context.Context) context.Context {
	lg := logger.With(
		zap.String("request_id", helpers.GenerateUUID()),
	)
	return context.WithValue(ctx, loggerKey{}, lg)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*zap.Logger); ok {
		return l
	}

	// Return a no-op logger if none is found
	return zap.NewNop()
}

type FiberContextWrapper struct {
	context fiber.Ctx
}

func NewFiberContextWrapper(c fiber.Ctx) pkg.AccessTokenContextWrapper {
	return &FiberContextWrapper{
		context: c,
	}
}

func (f *FiberContextWrapper) Context() context.Context {
	return f.context.Context()
}

func (f *FiberContextWrapper) GetToken() (string, error) {
	accessToken := f.context.Cookies(pkg.TOKEN_COOKIE, "")
	return accessToken, nil
}

func (f *FiberContextWrapper) DeleteToken() {
	f.context.ClearCookie(pkg.TOKEN_COOKIE)
}

func (f *FiberContextWrapper) SetToken(value string, durationInMinutes uint) {
	f.context.Cookie(&fiber.Cookie{
		Name:     pkg.TOKEN_COOKIE,
		Value:    value,
		Expires:  time.Now().Add(time.Minute * time.Duration(durationInMinutes)),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
}
