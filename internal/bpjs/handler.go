package bpjs

import (
	"context"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type service interface {
	GetPesertaByNIK(ctx context.Context, nik string) (*PesertaResponse, error)
}

type Handler struct {
	svc service
}

func NewHandler(svc service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetPesertaByNIK(c *echo.Context) error {
	nik := strings.TrimSpace(c.Param("nik"))

	peserta, err := h.svc.GetPesertaByNIK(c.Request().Context(), nik)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, peserta)
}
