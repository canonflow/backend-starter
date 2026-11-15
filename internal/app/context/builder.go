package context

import (
	"context"

	"github.com/canonflow/backend-starter/pkg/helpers"
	"go.uber.org/zap"
)

type loggerKey struct{}

func WithLogger(ctx context.Context) context.Context {
	lg := GetLogger().With(
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

func GetLocal[T any](ctx *FiberContextWrapper, key string) (T, bool) {
	var zero T
	raw := ctx.Local(key)
	if raw == nil {
		return zero, false
	}

	val, ok := raw.(T)
	if !ok {
		return zero, false
	}

	return val, true
}

func SetLocal[T any](ctx *FiberContextWrapper, key string, value T) {
	ctx.SetLocal(key, value)
}
