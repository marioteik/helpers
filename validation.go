package helpers

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator(validator *validator.Validate) *CustomValidator {
	return &CustomValidator{validator: validator}
}

func (cv *CustomValidator) Validate(data any) error {
	err := cv.validator.Struct(data)
	if err != nil {
		return err
	}

	return nil
}
