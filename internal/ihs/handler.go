package ihs

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc PatientService
}

func NewPatientHandler(svc PatientService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListPatients(c *echo.Context) error {
	var params GetPatientParams
	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	patients, total, err := h.svc.ListPatients(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, patients, params.Page, params.Limit, total)
}
