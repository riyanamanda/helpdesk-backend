package auth

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type handler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) handler {
	return handler{
		svc: svc,
	}
}

func (h *handler) Login(c *echo.Context) error {
	req, err := request.BindAndValidate[LoginRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	result, err := h.svc.Login(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, result)
}
