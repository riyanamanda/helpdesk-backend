package ticket_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
	ticket "github.com/riyanamanda/helpdesk-backend/internal/ticket"
)

type mockTicketRepository struct {
	mock.Mock
}

func (m *mockTicketRepository) GetAll(ctx context.Context, params ticket.GetTicketParams) ([]ticket.Ticket, int, error) {
	args := m.Called(ctx, params)

	var items []ticket.Ticket
	if v := args.Get(0); v != nil {
		items = v.([]ticket.Ticket)
	}

	return items, args.Int(1), args.Error(2)
}

func TestService_FetchAllTickets(t *testing.T) {
	now := time.Now().UTC()

	testCases := []struct {
		name      string
		params    *ticket.GetTicketParams
		setupMock func(*mockTicketRepository)
		assertFn  func(*testing.T, []ticket.TicketResponse, int, error)
	}{
		{
			name:   "success with default pagination",
			params: nil,
			setupMock: func(repo *mockTicketRepository) {
				items := []ticket.Ticket{
					{
						ID:          1,
						Title:       "Laptop issue",
						Description: "Cannot boot",
						CategoryID:  2,
						Status:      string(ticket.OPEN),
						Priority:    string(ticket.HIGH),
						CreatedBy:   uuid.New(),
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				}

				repo.On("GetAll", mock.Anything, ticket.GetTicketParams{Params: pagination.Params{Page: 1, Limit: 10}}).
					Return(items, 1, nil).Once()
			},
			assertFn: func(t *testing.T, result []ticket.TicketResponse, total int, err error) {
				require.NoError(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, "Laptop issue", result[0].Title)
				assert.Equal(t, 1, total)
			},
		},
		{
			name:   "repository error",
			params: &ticket.GetTicketParams{},
			setupMock: func(repo *mockTicketRepository) {
				repo.On("GetAll", mock.Anything, mock.Anything).
					Return(nil, 0, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result []ticket.TicketResponse, total int, err error) {
				require.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, 0, total)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockTicketRepository)
			t.Cleanup(func() { repo.AssertExpectations(t) })

			svc := ticket.NewTicketService(repo)
			tc.setupMock(repo)

			result, total, err := svc.FetchAllTickets(context.Background(), tc.params)
			tc.assertFn(t, result, total, err)
		})
	}
}
