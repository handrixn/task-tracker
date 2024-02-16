package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidatePayload(payload interface{}) map[string]string {
	validate := validator.New()

	if err := validate.Struct(payload); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			validationErrors[field] = fmt.Sprintf("%v is %v", field, err.Tag())
		}

		return validationErrors
	}

	return nil
}
