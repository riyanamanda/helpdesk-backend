package ticket

import "errors"

var (
	ErrTicketNotFound                = errors.New("ticket not found")
	ErrTicketResolutionAlreadyExists = errors.New("ticket resolution already exists")
)
