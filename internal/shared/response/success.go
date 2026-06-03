package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type SuccessResponse struct {
	Data any `json:"data,omitempty"`
}

func OK(c *echo.Context, data any) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Data: data,
	})
}

func Created(c *echo.Context, data any) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Data: data,
	})
}

func NoContent(c *echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
