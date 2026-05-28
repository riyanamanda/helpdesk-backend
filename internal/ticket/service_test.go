package ticket_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
	ticket "github.com/riyanamanda/helpdesk-backend/internal/ticket"
	ticketmocks "github.com/riyanamanda/helpdesk-backend/internal/ticket/mocks"
	apperror "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
	testingutil "github.com/riyanamanda/helpdesk-backend/internal/shared/testing"
	storagemocks "github.com/riyanamanda/helpdesk-backend/internal/storage/mocks"
)

func newTestService(repo ticket.TicketRepository, storage *storagemocks.Storage) ticket.TicketService {
	return ticket.NewTicketService(repo, storage, config.Storage{})
}

func TestService_FetchAllTickets(t *testing.T) {
	testCases := []struct {
		name      string
		params    *ticket.GetTicketParams
		setupMock func(*ticketmocks.TicketRepository)
		assertFn  func(*testing.T, []ticket.TicketResponse, int64, error)
	}{
		{
			name:   "success",
			params: &ticket.GetTicketParams{},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				items := []ticket.TicketProjection{
					{ID: 1, Title: "Printer broken", Status: "OPEN"},
					{ID: 2, Title: "VPN issue", Status: "IN_PROGRESS"},
				}
				repo.On("GetAll", mock.Anything, mock.Anything).Return(items, int64(2), nil).Once()
			},
			assertFn: func(t *testing.T, result []ticket.TicketResponse, total int64, err error) {
				require.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name:   "repository error",
			params: &ticket.GetTicketParams{},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("GetAll", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result []ticket.TicketResponse, total int64, err error) {
				require.Error(t, err)
				assert.Empty(t, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := ticketmocks.NewTicketRepository(t)
			storage := storagemocks.NewStorage(t)
			svc := newTestService(repo, storage)
			tc.setupMock(repo)

			result, total, err := svc.FetchAllTickets(context.Background(), tc.params)
			tc.assertFn(t, result, total, err)
		})
	}
}

func TestService_FindTicketByID(t *testing.T) {
	testCases := []struct {
		name      string
		id        int64
		setupMock func(*ticketmocks.TicketRepository)
		assertFn  func(*testing.T, ticket.TicketDetailResponse, error)
	}{
		{
			name: "success",
			id:   1,
			setupMock: func(repo *ticketmocks.TicketRepository) {
				proj := &ticket.TicketProjection{ID: 1, Title: "Printer broken", Status: "OPEN"}
				attachments := &[]ticket.TicketAttachmentProjection{}
				repo.On("GetByID", mock.Anything, int64(1)).Return(proj, nil).Once()
				repo.On("GetAttachmentsByTicketID", mock.Anything, int64(1)).Return(attachments, nil).Once()
			},
			assertFn: func(t *testing.T, result ticket.TicketDetailResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, int64(1), result.ID)
				assert.Equal(t, "Printer broken", result.Title)
			},
		},
		{
			name: "not found",
			id:   2,
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("GetByID", mock.Anything, int64(2)).Return(nil, ticket.ErrTicketNotFound).Once()
			},
			assertFn: func(t *testing.T, result ticket.TicketDetailResponse, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeNotFound, http.StatusNotFound, "ticket not found")
			},
		},
		{
			name: "repository error",
			id:   3,
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("GetByID", mock.Anything, int64(3)).Return(nil, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result ticket.TicketDetailResponse, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := ticketmocks.NewTicketRepository(t)
			storage := storagemocks.NewStorage(t)
			svc := newTestService(repo, storage)
			tc.setupMock(repo)

			result, err := svc.FindTicketByID(context.Background(), tc.id)
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_SetPriority(t *testing.T) {
	testCases := []struct {
		name      string
		id        int64
		req       ticket.TicketPriorityRequest
		setupMock func(*ticketmocks.TicketRepository)
		assertFn  func(*testing.T, error)
	}{
		{
			name: "success",
			id:   1,
			req:  ticket.TicketPriorityRequest{Priority: ticket.High},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("UpdatePriority", mock.Anything, int64(1), ticket.High).Return(nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "not found",
			id:   2,
			req:  ticket.TicketPriorityRequest{Priority: ticket.High},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("UpdatePriority", mock.Anything, int64(2), ticket.High).Return(ticket.ErrTicketNotFound).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeNotFound, http.StatusNotFound, "ticket not found")
			},
		},
		{
			name: "repository error",
			id:   3,
			req:  ticket.TicketPriorityRequest{Priority: ticket.High},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("UpdatePriority", mock.Anything, int64(3), ticket.High).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := ticketmocks.NewTicketRepository(t)
			storage := storagemocks.NewStorage(t)
			svc := newTestService(repo, storage)
			tc.setupMock(repo)

			err := svc.SetPriority(context.Background(), tc.id, tc.req)
			tc.assertFn(t, err)
		})
	}
}

func TestService_AssignTicket(t *testing.T) {
	priorityStr := string(ticket.High)

	testCases := []struct {
		name      string
		ticketID  int64
		req       ticket.TicketAssignRequest
		setupMock func(*ticketmocks.TicketRepository)
		assertFn  func(*testing.T, error)
	}{
		{
			name:     "success",
			ticketID: 1,
			req:      ticket.TicketAssignRequest{AssignedTo: uuid.New()},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				proj := &ticket.TicketProjection{ID: 1, Priority: &priorityStr}
				repo.On("GetByID", mock.Anything, int64(1)).Return(proj, nil).Once()
				repo.On("Assign", mock.Anything, int64(1), mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "priority not set",
			ticketID: 2,
			req:      ticket.TicketAssignRequest{AssignedTo: uuid.New()},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				proj := &ticket.TicketProjection{ID: 2, Priority: nil}
				repo.On("GetByID", mock.Anything, int64(2)).Return(proj, nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeBadRequest, http.StatusBadRequest, "Please set priority before assign a ticket")
			},
		},
		{
			name:     "repository error",
			ticketID: 3,
			req:      ticket.TicketAssignRequest{AssignedTo: uuid.New()},
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("GetByID", mock.Anything, int64(3)).Return(nil, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := ticketmocks.NewTicketRepository(t)
			storage := storagemocks.NewStorage(t)
			svc := newTestService(repo, storage)
			tc.setupMock(repo)

			err := svc.AssignTicket(context.Background(), tc.ticketID, tc.req)
			tc.assertFn(t, err)
		})
	}
}

func TestService_CloseTicket(t *testing.T) {
	testCases := []struct {
		name      string
		ticketID  int64
		setupMock func(*ticketmocks.TicketRepository)
		assertFn  func(*testing.T, error)
	}{
		{
			name:     "success",
			ticketID: 1,
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("CloseTicket", mock.Anything, int64(1), mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "not found",
			ticketID: 2,
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("CloseTicket", mock.Anything, int64(2), mock.AnythingOfType("uuid.UUID")).Return(ticket.ErrTicketNotFound).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				testingutil.AssertAppError(t, err, apperror.CodeNotFound, http.StatusNotFound, "ticket not found")
			},
		},
		{
			name:     "repository error",
			ticketID: 3,
			setupMock: func(repo *ticketmocks.TicketRepository) {
				repo.On("CloseTicket", mock.Anything, int64(3), mock.AnythingOfType("uuid.UUID")).Return(errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := ticketmocks.NewTicketRepository(t)
			storage := storagemocks.NewStorage(t)
			svc := newTestService(repo, storage)
			tc.setupMock(repo)

			err := svc.CloseTicket(context.Background(), tc.ticketID)
			tc.assertFn(t, err)
		})
	}
}

func TestService_RegisterTicket(t *testing.T) {
	testCases := []struct {
		name      string
		req       *ticket.TicketCreateRequest
		setupMock func(*ticketmocks.TicketRepository, *ticketmocks.TicketTx)
		assertFn  func(*testing.T, error)
	}{
		{
			name: "success without attachment",
			req:  &ticket.TicketCreateRequest{Title: "Issue", Description: "desc", CategoryID: 1, DivisionID: 1},
			setupMock: func(repo *ticketmocks.TicketRepository, tx *ticketmocks.TicketTx) {
				repo.On("Begin", mock.Anything).Return(tx, nil).Once()
				tx.On("Create", mock.Anything, mock.Anything).Return(int64(1), nil).Once()
				tx.On("Commit").Return(nil).Once()
				tx.On("Rollback").Return(nil).Maybe()
			},
			assertFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "begin transaction error",
			req:  &ticket.TicketCreateRequest{Title: "Issue", Description: "desc", CategoryID: 1, DivisionID: 1},
			setupMock: func(repo *ticketmocks.TicketRepository, tx *ticketmocks.TicketTx) {
				repo.On("Begin", mock.Anything).Return(nil, errors.New("tx error")).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "tx error")
			},
		},
		{
			name: "create error",
			req:  &ticket.TicketCreateRequest{Title: "Issue", Description: "desc", CategoryID: 1, DivisionID: 1},
			setupMock: func(repo *ticketmocks.TicketRepository, tx *ticketmocks.TicketTx) {
				repo.On("Begin", mock.Anything).Return(tx, nil).Once()
				tx.On("Create", mock.Anything, mock.Anything).Return(int64(0), errors.New("insert error")).Once()
				tx.On("Rollback").Return(nil).Once()
			},
			assertFn: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "insert error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := ticketmocks.NewTicketRepository(t)
			tx := ticketmocks.NewTicketTx(t)
			storage := storagemocks.NewStorage(t)
			svc := newTestService(repo, storage)
			tc.setupMock(repo, tx)

			err := svc.RegisterTicket(context.Background(), tc.req, nil)
			tc.assertFn(t, err)
		})
	}
}
