package category

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/response"
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
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	params := ListCategoriesParams{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.svc.GetCategories(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("internal server error"))
	}

	meta := &response.Meta{
		Limit:  params.Limit,
		Offset: params.Offset,
		Total:  result.Total,
	}

	return c.JSON(http.StatusOK,
		response.Success(toCategoryResponses(result.Data), meta),
	)
}
