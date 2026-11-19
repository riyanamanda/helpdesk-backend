package apperr

import (
	"errors"
)

type Kind int

const (
	KindInternal Kind = iota
	KindInvalid
	KindUnauthenticated
	KindForbidden
	KindNotFound
	KindConflict
	KindRateLimited
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
	Details map[string]any
}

func (e *Error) Error() string {
	return e.Message
}

func NotFound(resource string) *Error {
	return &Error{
		Kind:    KindNotFound,
		Code:    CodeNotFound,
		Message: resource + " not found",
	}
}

func AlreadyExists(resource string) *Error {
	return &Error{
		Kind:    KindConflict,
		Code:    CodeAlreadyExists,
		Message: resource + " already exists",
	}
}

func Conflict(message string) *Error {
	return &Error{
		Kind:    KindConflict,
		Code:    CodeConflict,
		Message: message,
	}
}

func BadRequest(message string) *Error {
	return &Error{
		Kind:    KindInvalid,
		Code:    CodeBadRequest,
		Message: message,
	}
}

func Forbidden(message string) *Error {
	return &Error{
		Kind:    KindForbidden,
		Code:    CodeForbidden,
		Message: message,
	}
}

func Unauthorized(code string, message string) *Error {
	return &Error{
		Kind:    KindUnauthenticated,
		Code:    code,
		Message: message,
	}
}

func Internal() *Error {
	return &Error{
		Kind:    KindInternal,
		Code:    CodeInternalError,
		Message: "internal server error",
	}
}

func RateLimited(message string, retryAfter int64) *Error {
	return &Error{
		Kind:    KindRateLimited,
		Code:    CodeRateLimited,
		Message: message,
		Details: map[string]any{
			"retry_after": retryAfter,
		},
	}
}

func InvalidResetToken(message string) *Error {
	return &Error{
		Kind:    KindInvalid,
		Code:    CodeInvalidResetToken,
		Message: message,
	}
}

func InvalidCredentials() *Error {
	return &Error{
		Kind:    KindInvalid,
		Code:    CodeInvalidCredentials,
		Message: "invalid email or password",
	}
}

func As(err error) *Error {
	if e, ok := errors.AsType[*Error](err); ok {
		return e
	}

	return Internal()
}
