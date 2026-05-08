package category

import (
	"net/http"

	"github.com/labstack/echo/v5"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
	sharedutils "github.com/riyanamanda/helpdesk-backend/internal/shared/utils"
)

type handler struct {
	svc CategoryService
}

func NewCategoryHandler(svc CategoryService) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) ListCategories(c *echo.Context) error {
	var params GetCategoryParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperror.BadRequest("invalid query params"))
	}

	page, limit, _ := params.Normalize()

	categories, total, err := h.svc.GetCategories(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, categories, page, limit, total)
}

func (h *handler) Create(c *echo.Context) error {
	req, err := request.BindAndValidate[CreateCategoryRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, category)
}

func (h *handler) GetByID(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, category)
}

func (h *handler) Update(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[UpdateCategoryRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.Update(c.Request().Context(), id, req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, category)
}

func (h *handler) Delete(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
