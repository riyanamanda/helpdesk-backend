package division

import "errors"

var (
	ErrDivisionAlreadyExists = errors.New("division already exists")
	ErrDivisionNotFound      = errors.New("division not found")
)
