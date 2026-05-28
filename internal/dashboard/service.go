package dashboard

import "context"

type DashboardService interface {
	GetSummary(ctx context.Context) (SummaryResponse, error)
	GetRecentTickets(ctx context.Context) ([]RecentTicket, error)
}

type service struct {
	repo DashboardRepository
}

func NewDashboardService(repo DashboardRepository) DashboardService {
	return &service{
		repo: repo,
	}
}

func (s *service) GetSummary(ctx context.Context) (SummaryResponse, error) {
	return s.repo.GetSummary(ctx)
}

func (s *service) GetRecentTickets(ctx context.Context) ([]RecentTicket, error) {
	return s.repo.GetRecentTickets(ctx)
}
