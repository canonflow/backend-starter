package app

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type loggerKey struct{}

func NewContext(ctx context.Context) context.Context {
	lg := logger.With(
		zap.String("request_id", uuid.NewString()),
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
