package profile

import "errors"

var (
	ErrProfileNotFound       = errors.New("profile not found")
	ErrGoogleIDAlreadyLinked = errors.New("google account already linked")
)
