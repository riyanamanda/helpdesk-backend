package dashboard

import "context"

type DashboardService interface {
	GetSummary(ctx context.Context) (SummaryResponse, error)
	GetRecentTickets(ctx context.Context) ([]RecentTicketResponse, error)
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
	summary, err := s.repo.GetSummary(ctx)
	if err != nil {
		return SummaryResponse{}, err
	}

	return toSummary(summary), nil
}

func (s *service) GetRecentTickets(ctx context.Context) ([]RecentTicketResponse, error) {
	recentTickets, err := s.repo.GetRecentTickets(ctx)
	if err != nil {
		return nil, err
	}

	return toRecentTickets(recentTickets), nil
}
