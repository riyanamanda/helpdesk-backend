package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func Error(c *echo.Context, err error) error {
	appErr := apperr.As(err)
	status := StatusOf(err)

	if status >= 500 {
		c.Set("internal_error", err)
	}

	return c.JSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}

func StatusOf(err error) int {
	switch apperr.As(err).Kind {
	case apperr.KindInvalid:
		return http.StatusBadRequest
	case apperr.KindUnauthenticated:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	case apperr.KindRateLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
