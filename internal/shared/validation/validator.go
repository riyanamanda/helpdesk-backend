package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func Parse(err error) *apperror.AppError {
	return apperror.Validation(
		apperror.ValidationErrors(err),
	)
}
