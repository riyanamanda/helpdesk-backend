package ihs

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

const maxNormLength = 32

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

func (h *Handler) GetPatientDetail(c *echo.Context) error {
	norm, err := httputil.ParseRequiredStringPathParam(c, "norm", "patient", maxNormLength)
	if err != nil {
		return response.Error(c, err)
	}

	patient, err := h.svc.GetPatientByNORM(c.Request().Context(), norm)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, patient)
}

func (h *Handler) UpdatePatientMethod(c *echo.Context) error {
	norm, err := httputil.ParseRequiredStringPathParam(c, "norm", "patient", maxNormLength)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdatePatientMethodByNORM(c.Request().Context(), norm); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
