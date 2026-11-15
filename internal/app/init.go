package app

import (
	appContext "github.com/canonflow/backend-starter/internal/app/context"
	goValidator "github.com/go-playground/validator/v10"
)

var validator *goValidator.Validate

func init() {
	appContext.CreateLogger()
	validator = newValidator()
}
