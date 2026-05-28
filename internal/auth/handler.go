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

func NewAuthHandler(svc AuthService) *handler {
	return &handler{
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

func (h *handler) LoginWithGoogle(c *echo.Context) error {
	req, err := request.BindAndValidate[GoogleLoginRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	result, err := h.svc.LoginWithGoogle(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, result)
}

func (h *handler) Me(c *echo.Context) error {
	user, err := h.svc.Me(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, user)
}
