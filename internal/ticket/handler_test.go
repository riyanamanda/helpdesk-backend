package ticket_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	ticket "github.com/riyanamanda/helpdesk-backend/internal/ticket"
)

type mockTicketService struct {
	mock.Mock
}

func (m *mockTicketService) FetchAllTickets(ctx context.Context, params *ticket.GetTicketParams) ([]ticket.TicketResponse, int, error) {
	args := m.Called(ctx, params)

	var items []ticket.TicketResponse
	if v := args.Get(0); v != nil {
		items = v.([]ticket.TicketResponse)
	}

	return items, args.Int(1), args.Error(2)
}

func TestHandler_ListTickets(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/tickets?page=1&limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		svc := new(mockTicketService)
		t.Cleanup(func() { svc.AssertExpectations(t) })

		h := ticket.NewTicketHandler(svc)

		svc.On("FetchAllTickets", mock.Anything, mock.MatchedBy(func(params *ticket.GetTicketParams) bool {
			return params != nil && params.Page == 1 && params.Limit == 10
		})).Return([]ticket.TicketResponse{{ID: 1, Title: "Printer issue"}}, 1, nil).Once()

		err := h.ListTickets(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Printer issue")
	})

	t.Run("invalid query params", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/tickets?page=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		svc := new(mockTicketService)
		t.Cleanup(func() { svc.AssertExpectations(t) })

		h := ticket.NewTicketHandler(svc)

		err := h.ListTickets(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid query params")
	})

	t.Run("service error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/tickets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		svc := new(mockTicketService)
		t.Cleanup(func() { svc.AssertExpectations(t) })

		h := ticket.NewTicketHandler(svc)

		svc.On("FetchAllTickets", mock.Anything, mock.Anything).
			Return(nil, 0, apperror.Internal("internal server error")).Once()

		err := h.ListTickets(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")
	})
}
