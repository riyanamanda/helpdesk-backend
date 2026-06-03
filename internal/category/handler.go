package category

import (
	"github.com/labstack/echo/v5"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
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
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	categories, total, err := h.svc.ListCategories(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, categories, params.Page, params.Limit, total)
}

func (h *Handler) ListCategoryOptions(c *echo.Context) error {
	categories, err := h.svc.ListOptions(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, categories)
}

func (h *Handler) CreateCategory(c *echo.Context) error {
	req, err := request.BindAndValidate[CategoryCreateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	category, err := h.svc.CreateCategory(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, category)
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

	return response.OK(c, category)
}

func (h *Handler) UpdateCategory(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := request.BindAndValidate[CategoryUpdateRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdateCategory(c.Request().Context(), id, req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) DeleteCategory(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "category")
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.DeleteCategory(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
