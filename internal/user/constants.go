package user

type UserRole string

const (
	ADMIN    UserRole = "ADMIN"
	EMPLOYEE UserRole = "EMPLOYEE"

	AssignableCacheKey = "user:assignable"
)
