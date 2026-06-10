package user_device

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc UserDeviceService
}

func NewUserDeviceHandler(svc UserDeviceService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterDevice(c *echo.Context) error {
	req, err := request.BindAndValidate[RegisterDeviceRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.RegisterDevice(c.Request().Context(), *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) UnregisterDevice(c *echo.Context) error {
	req, err := request.BindAndValidate[UnregisterDeviceRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UnregisterDevice(c.Request().Context(), *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
