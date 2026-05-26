package division

import (
	"net/http"

	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	sharedutils "github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
)

type handler struct {
	svc DivisionService
}

func NewDivisionHandler(svc DivisionService) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) ListDivisions(c *echo.Context) error {
	var params GetDivisionParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperror.BadRequest("invalid query params"))
	}

	divisions, total, err := h.svc.FetchAllDivisions(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, divisions, params.Page, params.Limit, total)
}

func (h *handler) CreateDivision(c *echo.Context) error {
	req, err := request.BindAndValidate[CreateDivisionRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	division, err := h.svc.RegisterDivision(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, division)
}

func (h *handler) GetDivision(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "division")
	if err != nil {
		return response.Error(c, err)
	}

	division, err := h.svc.FindDivisionByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, division)
}

func (h *handler) UpdateDivision(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "division")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[UpdateDivisionRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.EditDivision(c.Request().Context(), id, req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "division updated successfully")
}

func (h *handler) DeleteDivision(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "division")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.DeleteDivision(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
