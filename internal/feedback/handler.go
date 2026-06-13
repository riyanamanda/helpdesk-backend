package feedback

import (
	"github.com/labstack/echo/v5"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/httputil"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/response"
)

type Handler struct {
	svc FeedbackService
}

func NewFeedbackHandler(svc FeedbackService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListFeedbacks(c *echo.Context) error {
	var params GetFeedbackParams
	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	userID, ok := ctxkey.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return response.Error(c, apperr.Unauthorized(apperr.CodeUnauthorized, "unauthorized"))
	}
	params.CreatedByID = &userID

	feedbacks, total, err := h.svc.ListFeedbacks(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, feedbacks, params.Page, params.Limit, total)
}

func (h *Handler) ListAllFeedbacks(c *echo.Context) error {
	var params GetFeedbackParams
	if err := c.Bind(&params); err != nil {
		return response.Error(c, apperr.BadRequest("invalid query params"))
	}

	feedbacks, total, err := h.svc.ListFeedbacks(c.Request().Context(), &params)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Paginated(c, feedbacks, params.Page, params.Limit, total)
}

func (h *Handler) CreateFeedback(c *echo.Context) error {
	req, err := httputil.BindAndValidate[CreateFeedbackRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.CreateFeedback(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

func (h *Handler) GetFeedback(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "feedback")
	if err != nil {
		return response.Error(c, err)
	}

	feedback, err := h.svc.GetFeedback(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, feedback)
}

func (h *Handler) UpdateFeedbackStatus(c *echo.Context) error {
	id, err := httputil.ParsePositiveInt64PathParam(c, "id", "feedback")
	if err != nil {
		return response.Error(c, err)
	}

	req, err := httputil.BindAndValidate[UpdateFeedbackStatusRequest](c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.svc.UpdateFeedbackStatus(c.Request().Context(), id, *req); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}
