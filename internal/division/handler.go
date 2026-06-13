package division

import (
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc DivisionService
}

func NewDivisionHandler(svc DivisionService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) ListDivisions(c *echo.Context) error {
	var params GetDivisionParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	divisions, total, err := h.svc.ListDivisions(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, divisions, params.Page, params.Limit, total)
}

func (h *Handler) ListDivisionOptions(c *echo.Context) error {
	divisions, err := h.svc.ListOptions(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, divisions)
}

func (h *Handler) CreateDivision(c *echo.Context) error {
	req, err := httputil.BindAndValidate[DivisionCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.CreateDivision(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) GetDivision(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "division")
	if err != nil {
		return response.Error(c, err)
	}

	division, err := h.svc.GetDivision(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, division)
}

func (h *Handler) UpdateDivision(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "division")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := httputil.BindAndValidate[DivisionUpdateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdateDivision(c.Request().Context(), id, req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) DeleteDivision(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "division")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.DeleteDivision(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
