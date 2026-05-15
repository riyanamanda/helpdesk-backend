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

	categories, total, err := h.svc.FetchAllCategories(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, categories, params.Page, params.Limit, int64(total))
}

func (h *handler) CreateCategory(c *echo.Context) error {
	req, err := request.BindAndValidate[CreateCategoryRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.RegisterCategory(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, category)
}

func (h *handler) GetCategory(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.FindCategoryByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, category)
}

func (h *handler) UpdateCategory(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[UpdateCategoryRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.EditCategory(c.Request().Context(), id, req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "category updated successfully")
}

func (h *handler) DeleteCategory(c *echo.Context) error {
	id, err := sharedutils.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.DeleteCategory(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
