package profile

import "github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"

var (
	ErrProfileNotFound       = apperr.NotFound("profile")
	ErrGoogleIDAlreadyLinked = apperr.AlreadyExists("google account")
)
