package app

import (
	goValidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	logger    *zap.Logger
	validator *goValidator.Validate
)

func init() {
	logger = createLogger()
	validator = newValidator()

	defer logger.Sync()
}
