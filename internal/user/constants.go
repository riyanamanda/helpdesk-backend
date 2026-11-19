package user

import "time"

const (
	AssignableCacheKey = "user:assignable"
)

const (
	LoginMethodPassword = "PASSWORD"
	LoginMethodGoogle   = "GOOGLE"
)

const (
	SetPasswordTokenTTL   = 72 * time.Hour
	SetPasswordExpiryText = "72 jam"
)
