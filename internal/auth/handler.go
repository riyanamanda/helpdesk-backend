package auth

import (
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Login(c *echo.Context) error {
	req, err := request.BindAndValidate[LoginRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	result, err := h.svc.Login(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, result)
}

func (h *Handler) LoginWithGoogle(c *echo.Context) error {
	req, err := request.BindAndValidate[GoogleLoginRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	result, err := h.svc.LoginWithGoogle(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, result)
}

func (h *Handler) Logout(c *echo.Context) error {
	if err := h.svc.Logout(c.Request().Context()); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) Me(c *echo.Context) error {
	user, err := h.svc.Me(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, user)
}
