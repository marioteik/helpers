package helpers

import "github.com/go-playground/validator/v10"

func ValidateStruct(data any) error {
	validate := validator.New()
	err := validate.Struct(data)

	if err != nil {
		return err
	}

	return nil
}
