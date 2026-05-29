package response

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
)

type Response[T any] struct {
	Data    *T      `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	Meta    Meta    `json:"meta"`
}

type ErrorResponse struct {
	Error ErrorInfo `json:"error"`
	Meta  Meta      `json:"meta"`
}

type ErrorInfo struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Meta struct {
	Timestamp  string      `json:"timestamp"`
	RequestID  string      `json:"request_id"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}

func getRequestID(c *echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func calculateTotalPage(total int64, limit int) int64 {
	if total == 0 {
		return 0
	}

	return (total + int64(limit) - 1) / int64(limit)
}

func buildMeta(c *echo.Context, pagination *Pagination) *Meta {
	return &Meta{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		RequestID:  getRequestID(c),
		Pagination: pagination,
	}
}

func WithPagination[T any](c *echo.Context, statusCode int, data T, page, limit int, total int64) error {
	pagination := &Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: calculateTotalPage(total, limit),
	}

	return c.JSON(statusCode, Response[T]{
		Data: &data,
		Meta: *buildMeta(c, pagination),
	})
}

func Message(c *echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Response[any]{
		Message: &message,
		Meta:    *buildMeta(c, nil),
	})
}

func Success[T any](c *echo.Context, statusCode int, data T) error {
	return c.JSON(statusCode, Response[T]{
		Data: &data,
		Meta: *buildMeta(c, nil),
	})
}

func Error(c *echo.Context, err error) error {
	appErr := apperror.As(err)
	if appErr.StatusCode >= 500 {
		slog.Error(
			"internal error",
			"request_id", getRequestID(c),
			"error", err,
		)
	}

	var details any
	if appErr.Details != nil {
		details = appErr.Details
	}

	errorInfo := ErrorInfo{
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: details,
	}

	return c.JSON(appErr.StatusCode, ErrorResponse{
		Error: errorInfo,
		Meta:  *buildMeta(c, nil),
	})
}
