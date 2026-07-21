package config

import "github.com/go-playground/validator/v10"

type IValidatorConfig interface {
    GetValidator() *validator.Validate
	RegisterValidation(tag string, validation func(fl validator.FieldLevel) bool) error
}