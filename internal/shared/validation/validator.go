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

	if err := v.RegisterValidation("notcommon", notCommonPassword); err != nil {
		panic(err)
	}

	return &CustomValidator{validator: v}
}

var commonPasswords = map[string]struct{}{
	"password":    {},
	"password1":   {},
	"password12":  {},
	"password123": {},
	"passw0rd":    {},
	"12345678":    {},
	"123456789":   {},
	"1234567890":  {},
	"87654321":    {},
	"11111111":    {},
	"00000000":    {},
	"qwerty123":   {},
	"qwertyuiop":  {},
	"abcd1234":    {},
	"iloveyou":    {},
	"admin123":    {},
	"helpdesk":    {},
	"helpdesk123": {},
	"rahasia123":  {},
}

func notCommonPassword(fl validator.FieldLevel) bool {
	_, common := commonPasswords[strings.ToLower(strings.TrimSpace(fl.Field().String()))]
	return !common
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func Parse(err error) *apperr.Error {
	return apperr.Validation(apperr.ValidationErrors(err))
}
