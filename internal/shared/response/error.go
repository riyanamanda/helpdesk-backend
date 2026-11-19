package response

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

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

func HTTPErrorHandler(c *echo.Context, err error) {
	if resp, _ := echo.UnwrapResponse(c.Response()); resp != nil && resp.Committed {
		return
	}

	status := StatusOf(err)

	var writeErr error
	if c.Request().Method == http.MethodHead {
		writeErr = c.NoContent(status)
	} else {
		writeErr = c.JSON(status, ErrorResponse{Error: bodyOf(err, status)})
	}

	if writeErr != nil {
		slog.WarnContext(c.Request().Context(), "write error response", "error", writeErr)
	}
}

func StatusOf(err error) int {
	if appErr, ok := errors.AsType[*apperr.Error](err); ok {
		return statusOfKind(appErr.Kind)
	}

	if status := echo.StatusCode(err); status != 0 {
		return status
	}

	return http.StatusInternalServerError
}

func statusOfKind(kind apperr.Kind) int {
	switch kind {
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

func bodyOf(err error, status int) ErrorBody {
	if appErr, ok := errors.AsType[*apperr.Error](err); ok {
		return ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		}
	}

	switch {
	case status >= http.StatusInternalServerError:
		return ErrorBody{Code: apperr.CodeInternalError, Message: "internal server error"}
	case status == http.StatusNotFound:
		return ErrorBody{Code: apperr.CodeNotFound, Message: "route not found"}
	case status == http.StatusMethodNotAllowed:
		return ErrorBody{Code: apperr.CodeMethodNotAllowed, Message: "method not allowed"}
	default:
		return ErrorBody{Code: apperr.CodeBadRequest, Message: strings.ToLower(http.StatusText(status))}
	}
}
