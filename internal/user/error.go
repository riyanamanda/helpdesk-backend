package user

import "github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"

var (
	ErrUserAlreadyExists = apperr.AlreadyExists("user")
	ErrUserNotFound      = apperr.NotFound("user")
)
