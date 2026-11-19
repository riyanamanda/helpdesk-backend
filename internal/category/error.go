package category

import "github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"

var (
	ErrCategoryAlreadyExists = apperr.AlreadyExists("category")
	ErrCategoryNotFound      = apperr.NotFound("category")
)
