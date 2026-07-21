package config

import "github.com/go-playground/validator/v10"

type ValidatorConfigStruct struct {
	validator *validator.Validate
}

func (v *ValidatorConfigStruct) GetValidator() *validator.Validate {
	if v.validator == nil {
		v.validator = validator.New()
	}

	return v.validator
}

func (v *ValidatorConfigStruct) RegisterValidation(tag string, validation func(fl validator.FieldLevel) bool) error {
	err := v.validator.RegisterValidation(tag, validation)

	return err
}

func NewValidatorConfigStruct() *ValidatorConfigStruct {
	return &ValidatorConfigStruct{
		validator: validator.New(),
	}
}