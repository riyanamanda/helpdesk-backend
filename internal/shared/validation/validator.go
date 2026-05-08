package validation

import (
	"github.com/go-playground/validator/v10"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func Parse(err error) *apperror.AppError {
	return apperror.Validation(
		apperror.ValidationErrors(err),
	)
}
