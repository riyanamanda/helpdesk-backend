package request

import (
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/validation"
)

func BindAndValidate[T any](c *echo.Context) (*T, error) {

	var req T

	if err := c.Bind(&req); err != nil {

		return nil, apperror.BadRequest("invalid request body")

	}

	if err := c.Validate(&req); err != nil {

		return nil, validation.Parse(err)

	}

	return &req, nil

}
