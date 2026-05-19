package apperror

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	CODE_NOT_FOUND        = "NOT_FOUND"
	CODE_ALREADY_EXISTS   = "ALREADY_EXISTS"
	CODE_VALIDATION_ERROR = "VALIDATION_ERROR"
	CODE_INTERNAL_ERROR   = "INTERNAL_SERVER_ERROR"
	CODE_BAD_REQUEST      = "BAD_REQUEST"
	CODE_FORBIDDEN        = "FORBIDDEN"
	CODE_UNAUTHORIZED = "UNAUTHORIZED"
	CODE_INVALID_TOKEN = "INVALID_TOKEN"
	CODE_TOKEN_EXPIRED = "TOKEN_EXPIRED"
	CODE_MISSING_TOKEN = "MISSING_TOKEN"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInternal      = errors.New("internal server error")
	ErrBadRequest    = errors.New("bad request")
	ErrForbidden     = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")
)

type AppError struct {
	Err        error
	Code       string
	Message    string
	StatusCode int
	Details    map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func NotFound(resource string) *AppError {
	return &AppError{
		Err:        ErrNotFound,
		Code:       CODE_NOT_FOUND,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func AlreadyExists(resource string) *AppError {
	return &AppError{
		Err:        ErrAlreadyExists,
		Code:       CODE_ALREADY_EXISTS,
		Message:    fmt.Sprintf("%s already exists", resource),
		StatusCode: http.StatusConflict,
	}
}

func Unauthorized(code string, message string) *AppError {
	return &AppError{
		Err: ErrUnauthorized,
		Code: code,
		Message: message,
		StatusCode: http.StatusUnauthorized,
	}
}

func Internal(message string) *AppError {
	return &AppError{
		Err:        ErrInternal,
		Code:       CODE_INTERNAL_ERROR,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Err:        ErrBadRequest,
		Code:       CODE_BAD_REQUEST,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Err:        ErrForbidden,
		Code:       CODE_FORBIDDEN,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func Validation(details map[string]any) *AppError {
	return &AppError{
		Err:        ErrBadRequest,
		Code:       CODE_VALIDATION_ERROR,
		Message:    "validation failed",
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

func ValidationErrors(err error) map[string]any {
	errors := map[string]any{}

	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())

		switch err.Tag() {

		case "required":
			errors[field] = field + " is required"
		case "min":
			errors[field] = field + " minimum length is " + err.Param()
		case "max":
			errors[field] = field + " maximum length is " + err.Param()
		case "email":
			errors[field] = "invalid email format"
		default:
			errors[field] = "invalid value"
		}
	}

	return errors
}

func As(err error) *AppError {
	var appErr *AppError

	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal("internal server error")
}
