package apperr

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Validation(details map[string]any) *Error {
	return &Error{
		Code:    CodeValidationError,
		Message: "validation failed",
		Status:  http.StatusBadRequest,
		Details: details,
	}
}

func ValidationErrors(err error) map[string]any {
	result := map[string]any{}

	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {
		return result
	}

	for _, e := range ve {
		result[e.Field()] = map[string]any{
			"code":  e.Tag(),
			"param": e.Param(),
		}
	}

	return result
}
