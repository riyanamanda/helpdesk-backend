package dashboard_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dashboard "github.com/riyanamanda/helpdesk-backend/internal/dashboard"
	dashboardmocks "github.com/riyanamanda/helpdesk-backend/internal/dashboard/mocks"
)

func TestService_GetSummary(t *testing.T) {
	testCases := []struct {
		name      string
		setupMock func(*dashboardmocks.DashboardRepository)
		assertFn  func(*testing.T, dashboard.SummaryResponse, error)
	}{
		{
			name: "success",
			setupMock: func(repo *dashboardmocks.DashboardRepository) {
				result := dashboard.SummaryResponse{
					ByStatus: dashboard.TicketStatusStats{
						Open:       3,
						InProgress: 2,
						Resolved:   1,
						Closed:     0,
						Total:      6,
					},
					ByPriority: dashboard.TicketPriorityStats{
						Low:    1,
						Medium: 2,
						High:   1,
						Urgent: 1,
					},
				}
				repo.On("GetSummary", mock.Anything).Return(result, nil).Once()
			},
			assertFn: func(t *testing.T, result dashboard.SummaryResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, int64(6), result.ByStatus.Total)
				assert.Equal(t, int64(3), result.ByStatus.Open)
				assert.Equal(t, int64(2), result.ByPriority.Medium)
			},
		},
		{
			name: "repository error",
			setupMock: func(repo *dashboardmocks.DashboardRepository) {
				repo.On("GetSummary", mock.Anything).Return(dashboard.SummaryResponse{}, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result dashboard.SummaryResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, dashboard.SummaryResponse{}, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := dashboardmocks.NewDashboardRepository(t)
			svc := dashboard.NewDashboardService(repo)
			tc.setupMock(repo)

			result, err := svc.GetSummary(context.Background())
			tc.assertFn(t, result, err)
		})
	}
}

func TestService_GetRecentTickets(t *testing.T) {
	testCases := []struct {
		name      string
		setupMock func(*dashboardmocks.DashboardRepository)
		assertFn  func(*testing.T, []dashboard.RecentTicket, error)
	}{
		{
			name: "success",
			setupMock: func(repo *dashboardmocks.DashboardRepository) {
				priority := "HIGH"
				assignee := "Bob"
				items := []dashboard.RecentTicket{
					{ID: 1, Title: "Printer broken", Status: "OPEN", Priority: &priority, CreatedBy: "Alice", AssignedTo: &assignee, CreatedAt: "2026-05-01T00:00:00Z"},
					{ID: 2, Title: "VPN issue", Status: "IN_PROGRESS", Priority: nil, CreatedBy: "Charlie", AssignedTo: nil, CreatedAt: "2026-05-02T00:00:00Z"},
				}
				repo.On("GetRecentTickets", mock.Anything).Return(items, nil).Once()
			},
			assertFn: func(t *testing.T, result []dashboard.RecentTicket, err error) {
				require.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, int64(1), result[0].ID)
				assert.Equal(t, "Printer broken", result[0].Title)
				assert.Nil(t, result[1].Priority)
			},
		},
		{
			name: "repository error",
			setupMock: func(repo *dashboardmocks.DashboardRepository) {
				repo.On("GetRecentTickets", mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			assertFn: func(t *testing.T, result []dashboard.RecentTicket, err error) {
				require.Error(t, err)
				assert.Nil(t, result)
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := dashboardmocks.NewDashboardRepository(t)
			svc := dashboard.NewDashboardService(repo)
			tc.setupMock(repo)

			result, err := svc.GetRecentTickets(context.Background())
			tc.assertFn(t, result, err)
		})
	}
}
