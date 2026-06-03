package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
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

func Parse(err error) *apperr.Error {
	return apperr.Validation(apperr.ValidationErrors(err))
}
