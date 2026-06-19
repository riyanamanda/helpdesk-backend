package antrian

import (
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc AntrianService
}

func NewAntrianHandler(svc AntrianService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListAntrian(c *echo.Context) error {
	var params GetAntrianParams
	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	antrian, total, err := h.svc.ListAntrian(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, antrian, params.Page, params.Limit, total)
}

func (h *Handler) CheckIn(c *echo.Context) error {
	kodeBooking, err := strconv.ParseInt(c.Param("kode_booking"), 10, 64)
	if err != nil {
		return response.Error(c, apperr.BadRequest("invalid kode_booking"))
	}

	if err := h.svc.CheckIn(c.Request().Context(), kodeBooking); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
