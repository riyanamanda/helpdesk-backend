package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
)

type repository interface {
	GetSummary(ctx context.Context) (SummaryProjection, error)
	GetMonthlyTrend(ctx context.Context, year int) ([]MonthlyTrendProjection, error)
	GetTicketsByCategory(ctx context.Context) ([]CategoryTicketsProjection, error)
	GetAgentWorkload(ctx context.Context) ([]AgentWorkloadProjection, error)
}

type Service struct {
	repo  repository
	cache cache.Cache
}

func NewService(repo repository, cache cache.Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) GetSummary(ctx context.Context) (*SummaryResponse, error) {
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

func (s *Service) GetMonthlyTrend(ctx context.Context, year int) ([]MonthlyTrendResponse, error) {
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

func (s *Service) GetTicketsByCategory(ctx context.Context) ([]CategoryTicketsResponse, error) {
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

func (s *Service) GetAgentWorkload(ctx context.Context) ([]AgentWorkloadResponse, error) {
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
