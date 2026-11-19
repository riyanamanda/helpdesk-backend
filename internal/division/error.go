package division

import "github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"

var (
	ErrDivisionAlreadyExists = apperr.AlreadyExists("division")
	ErrDivisionNotFound      = apperr.NotFound("division")
)
