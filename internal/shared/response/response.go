package response

import (
	"time"

	"github.com/labstack/echo/v5"
	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type Response[T any] struct {
	Data    T       `json:"data"`
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
}

type Meta struct {
	Timestamp  string      `json:"timestamp"`
	RequestID  string      `json:"request_id"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func getRequestID(c *echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func calculateTotalPage(total, limit int) int {
	if total == 0 {
		return 0
	}

	return (total + limit - 1) / limit
}

func buildMeta(c *echo.Context, pagination *Pagination) *Meta {
	return &Meta{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		RequestID:  getRequestID(c),
		Pagination: pagination,
	}
}

func WithPagination[T any](c *echo.Context, statusCode int, data T, page, limit, total int) error {
	pagination := &Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: calculateTotalPage(total, limit),
	}
	return c.JSON(statusCode, Response[T]{
		Data: data,
		Meta: *buildMeta(c, pagination),
	})
}

func Error(c *echo.Context, err error) error {
	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		appErr = apperrors.Internal("internal server error")
	}

	errorInfo := ErrorInfo{
		Code:    appErr.Code,
		Message: appErr.Message,
	}

	return c.JSON(appErr.StatusCode, ErrorResponse{
		Error: errorInfo,
		Meta:  *buildMeta(c, nil),
	})
}
