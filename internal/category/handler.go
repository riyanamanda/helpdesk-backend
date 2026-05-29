package category

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/request"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc CategoryService
}

func NewCategoryHandler(svc CategoryService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) ListCategories(c *echo.Context) error {
	var params GetCategoryParams

	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperror.BadRequest("invalid query params"))
	}

	categories, total, err := h.svc.ListCategories(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.WithPagination(c, http.StatusOK, categories, params.Page, params.Limit, total)
}

func (h *Handler) CreateCategory(c *echo.Context) error {
	req, err := request.BindAndValidate[CreateCategoryRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.CreateCategory(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, category)
}

func (h *Handler) GetCategory(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.GetCategory(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, category)
}

func (h *Handler) UpdateCategory(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[UpdateCategoryRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdateCategory(c.Request().Context(), id, req); err != nil {
		return response.Error(c, err)
	}

	return response.Message(c, http.StatusOK, "category updated successfully")
}

func (h *Handler) DeleteCategory(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.DeleteCategory(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
