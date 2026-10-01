package auth

import (
	"context"

	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type service interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	LoginWithGoogle(ctx context.Context, req *GoogleLoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context) error
	Me(ctx context.Context) (*CurrentUserResponse, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}

type Handler struct {
	svc service
}

func NewHandler(svc service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Login(c *echo.Context) error {
	req, err := httputil.BindAndValidate[LoginRequest](c)
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
	req, err := httputil.BindAndValidate[GoogleLoginRequest](c)
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

func (h *Handler) ForgotPassword(c *echo.Context) error {
	req, err := httputil.BindAndValidate[ForgotPasswordRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.ForgotPassword(c.Request().Context(), *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) ResetPassword(c *echo.Context) error {
	req, err := httputil.BindAndValidate[ResetPasswordRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.ResetPassword(c.Request().Context(), *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
