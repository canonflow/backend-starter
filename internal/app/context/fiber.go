package context

import (
	"context"
	"strings"
	"time"

	"github.com/canonflow/backend-starter/internal/config"
	"github.com/gofiber/fiber/v3"
)

var tokenCookie = "Authorization"

type FiberContextWrapper struct {
	context fiber.Ctx
}

func NewFiberContextWrapper(c fiber.Ctx) *FiberContextWrapper {
	return &FiberContextWrapper{
		context: c,
	}
}

func (f *FiberContextWrapper) Context() context.Context {
	return f.context.Context()
}

func (f *FiberContextWrapper) GetToken() (string, error) {
	accessToken := f.context.Cookies(tokenCookie, "")
	return accessToken, nil
}

func (f *FiberContextWrapper) DeleteToken() {
	f.context.Cookie(&fiber.Cookie{
		Name:     tokenCookie,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: config.Get[bool](config.JWTHTTPOnly),
		Secure:   config.Get[bool](config.JWTSecure),
		SameSite: strings.ToTitle(config.Get[string](config.JWTSameSite)),
	})
}

func (f *FiberContextWrapper) SetToken(value string, duratioInMinutes uint) {
	f.context.Cookie(&fiber.Cookie{
		Name:     tokenCookie,
		Value:    value,
		Expires:  time.Now().Add(time.Minute * time.Duration(duratioInMinutes)),
		HTTPOnly: config.Get[bool](config.JWTHTTPOnly),
		Secure:   config.Get[bool](config.JWTSecure),
		SameSite: strings.ToTitle(config.Get[string](config.JWTSameSite)),
	})
}

func (f *FiberContextWrapper) SetLocal(key string, value any) {
	f.context.Locals(key, value)
}

func (f *FiberContextWrapper) Local(key string) any {
	return f.context.Locals(key)
}
