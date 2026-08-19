package app

import (
	"github.com/canonflow/backend-starter/internal/core"
	goValidator "github.com/go-playground/validator/v10"
)

var validator *goValidator.Validate

func init() {
	core.CreateLogger()
	validator = newValidator()
}
