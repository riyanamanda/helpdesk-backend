package apperr

import (
	"errors"
	"net/http"
)

type Error struct {
	Code    string
	Message string
	Status  int
	Details map[string]any
}

func (e *Error) Error() string {
	return e.Message
}

func NotFound(resource string) *Error {
	return &Error{
		Code:    CodeNotFound,
		Message: resource + " not found",
		Status:  http.StatusNotFound,
	}
}

func AlreadyExists(resource string) *Error {
	return &Error{
		Code:    CodeAlreadyExists,
		Message: resource + " already exists",
		Status:  http.StatusConflict,
	}
}

func BadRequest(message string) *Error {
	return &Error{
		Code:    CodeBadRequest,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func Forbidden(message string) *Error {
	return &Error{
		Code:    CodeForbidden,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func Unauthorized(code string, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func Internal() *Error {
	return &Error{
		Code:    CodeInternalError,
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}
}

func As(err error) *Error {
	if e, ok := errors.AsType[*Error](err); ok {
		return e
	}

	return Internal()
}
