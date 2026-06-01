package dashboard

import (
	"context"
	"encoding/json"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

type DashboardService interface {
	GetSummary(ctx context.Context) (SummaryResponse, error)
	GetRecentTickets(ctx context.Context) ([]RecentTicketResponse, error)
}

type service struct {
	repo  DashboardRepository
	cache cache.Cache
}

func NewDashboardService(repo DashboardRepository, cache cache.Cache) DashboardService {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) GetSummary(ctx context.Context) (SummaryResponse, error) {
	cached, err := s.cache.Get(ctx, SummaryCacheKey)
	if err == nil {
		var summary SummaryResponse

		if err := json.Unmarshal([]byte(cached), &summary); err == nil {
			return summary, nil
		}
	}

	projection, err := s.repo.GetSummary(ctx)
	if err != nil {
		return SummaryResponse{}, err
	}

	summary := toSummary(projection)

	data, err := json.Marshal(summary)
	if err == nil {
		_ = s.cache.Set(ctx, SummaryCacheKey, string(data), 30*time.Second)
	}

	return summary, nil
}

func (s *service) GetRecentTickets(ctx context.Context) ([]RecentTicketResponse, error) {
	cached, err := s.cache.Get(ctx, RecentTicketsCacheKey)
	if err == nil {
		var tickets []RecentTicketResponse

		if err := json.Unmarshal([]byte(cached), &tickets); err == nil {
			return tickets, nil
		}
	}

	recentTickets, err := s.repo.GetRecentTickets(ctx)
	if err != nil {
		return nil, err
	}

	tickets := toRecentTickets(recentTickets)

	data, err := json.Marshal(tickets)
	if err == nil {
		_ = s.cache.Set(ctx, RecentTicketsCacheKey, string(data), 30*time.Second)
	}

	return tickets, nil
}
