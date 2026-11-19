package ihs

import "github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"

var (
	ErrPatientNotFound    = apperr.NotFound("patient")
	ErrPatientNotEligible = apperr.Conflict("patient has already been submitted or is currently being processed")
)
