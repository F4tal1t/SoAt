package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Friendships(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				fmt.Printf("Validation failed for field %s: %s\n", e.Field(), e.Tag())
			}
		}
		return err
	}
	return nil
}
