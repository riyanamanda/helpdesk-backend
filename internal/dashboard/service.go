package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

type DashboardService interface {
	GetSummary(ctx context.Context) (*SummaryResponse, error)
	GetRecentTickets(ctx context.Context) ([]RecentTicketResponse, error)
	GetMonthlyTrend(ctx context.Context, year int) ([]MonthlyTrendResponse, error)
	GetTicketsByCategory(ctx context.Context) ([]CategoryTicketsResponse, error)
	GetAgentWorkload(ctx context.Context) ([]AgentWorkloadResponse, error)
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

func (s *service) GetSummary(ctx context.Context) (*SummaryResponse, error) {
	cached, err := s.cache.Get(ctx, SummaryCacheKey)
	if err == nil {
		var summary SummaryResponse
		if err := json.Unmarshal([]byte(cached), &summary); err == nil {
			return &summary, nil
		}
	}

	projection, err := s.repo.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	summary := toSummary(projection)

	if data, err := json.Marshal(summary); err == nil {
		_ = s.cache.Set(ctx, SummaryCacheKey, string(data), 30*time.Second)
	}

	return &summary, nil
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

	if data, err := json.Marshal(tickets); err == nil {
		_ = s.cache.Set(ctx, RecentTicketsCacheKey, string(data), 30*time.Second)
	}

	return tickets, nil
}

func (s *service) GetMonthlyTrend(ctx context.Context, year int) ([]MonthlyTrendResponse, error) {
	cacheKey := fmt.Sprintf(MonthlyTrendCacheKey, year)

	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var trend []MonthlyTrendResponse
		if err := json.Unmarshal([]byte(cached), &trend); err == nil {
			return trend, nil
		}
	}

	rows, err := s.repo.GetMonthlyTrend(ctx, year)
	if err != nil {
		return nil, err
	}

	trend := toMonthlyTrend(rows)

	if data, err := json.Marshal(trend); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return trend, nil
}

func (s *service) GetTicketsByCategory(ctx context.Context) ([]CategoryTicketsResponse, error) {
	cached, err := s.cache.Get(ctx, CategoryTicketsCacheKey)
	if err == nil {
		var categories []CategoryTicketsResponse
		if err := json.Unmarshal([]byte(cached), &categories); err == nil {
			return categories, nil
		}
	}

	rows, err := s.repo.GetTicketsByCategory(ctx)
	if err != nil {
		return nil, err
	}

	categories := toCategoryTickets(rows)

	if data, err := json.Marshal(categories); err == nil {
		_ = s.cache.Set(ctx, CategoryTicketsCacheKey, string(data), 30*time.Second)
	}

	return categories, nil
}

func (s *service) GetAgentWorkload(ctx context.Context) ([]AgentWorkloadResponse, error) {
	cached, err := s.cache.Get(ctx, AgentWorkloadCacheKey)
	if err == nil {
		var workload []AgentWorkloadResponse
		if err := json.Unmarshal([]byte(cached), &workload); err == nil {
			return workload, nil
		}
	}

	rows, err := s.repo.GetAgentWorkload(ctx)
	if err != nil {
		return nil, err
	}

	workload := toAgentWorkload(rows)

	if data, err := json.Marshal(workload); err == nil {
		_ = s.cache.Set(ctx, AgentWorkloadCacheKey, string(data), 30*time.Second)
	}

	return workload, nil
}
