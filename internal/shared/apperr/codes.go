package apperr

const (
	CodeNotFound        = "NOT_FOUND"
	CodeAlreadyExists   = "ALREADY_EXISTS"
	CodeConflict        = "CONFLICT"
	CodeValidationError = "VALIDATION_ERROR"
	CodeInternalError   = "INTERNAL_SERVER_ERROR"
	CodeBadRequest      = "BAD_REQUEST"
	CodeForbidden       = "FORBIDDEN"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeRateLimited     = "RATE_LIMITED"

	CodeMethodNotAllowed = "METHOD_NOT_ALLOWED"

	CodeInvalidToken = "INVALID_TOKEN"
	CodeTokenExpired = "TOKEN_EXPIRED"
	CodeMissingToken = "MISSING_TOKEN"

	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeInvalidResetToken  = "INVALID_RESET_TOKEN"
)
