package app

import goValidator "github.com/go-playground/validator/v10"

func newValidator() *goValidator.Validate {
	return goValidator.New()
}

func GetValidator() *goValidator.Validate {
	return validator
}
