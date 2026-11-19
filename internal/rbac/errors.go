package rbac

import "github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"

var (
	ErrRoleNotFound         = apperr.NotFound("role")
	ErrInvalidPermissionIDs = apperr.BadRequest("one or more permission IDs are invalid")
)
