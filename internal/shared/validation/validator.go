package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	govalidator "github.com/go-playground/validator/v10"
	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type EchoValidator struct {
	validator *govalidator.Validate
}

func New() *EchoValidator {
	v := govalidator.New(govalidator.WithRequiredStructEnabled())
	_ = v.RegisterValidation("notblank", func(fl govalidator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	return &EchoValidator{validator: v}
}

func (v *EchoValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		var validationErrs govalidator.ValidationErrors
		if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
			fieldErr := validationErrs[0]
			fieldName := jsonFieldName(i, fieldErr.StructField())
			if fieldName == "" {
				fieldName = strings.ToLower(fieldErr.StructField())
			}
			return apperrors.Validation(fmt.Sprintf("%s is invalid (%s)", fieldName, fieldErr.Tag()))
		}
		return apperrors.Validation("invalid request payload")
	}

	return nil
}

func jsonFieldName(i interface{}, structField string) string {
	t := reflect.TypeOf(i)
	if t == nil {
		return ""
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}

	field, ok := t.FieldByName(structField)
	if !ok {
		return ""
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return ""
	}

	name := strings.Split(jsonTag, ",")[0]
	if name == "" {
		return ""
	}

	return name
}
