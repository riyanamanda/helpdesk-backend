package division

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

type repository interface {
	GetAll(ctx context.Context, params GetDivisionParams) ([]Division, int64, error)
	GetOptions(ctx context.Context) ([]DivisionOptionProjection, error)
	Create(ctx context.Context, division *Division) error
	GetByID(ctx context.Context, id int64) (*Division, error)
	Update(ctx context.Context, id int64, division *Division) error
	Delete(ctx context.Context, id int64) error
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

func (s *Service) ListDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int64, error) {
	if params == nil {
		params = &GetDivisionParams{}
	}
	params.Normalize()

	divisions, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return nil, 0, err
	}

	return toDivisionResponses(divisions), total, nil
}

func (s *Service) ListOptions(ctx context.Context) ([]DivisionOptionResponse, error) {
	cached, err := s.cache.Get(ctx, DivisionOptionsCacheKey)
	if err == nil {
		var divisions []DivisionOptionResponse

		if err := json.Unmarshal([]byte(cached), &divisions); err == nil {
			return divisions, nil
		}
	}

	projection, err := s.repo.GetOptions(ctx)
	if err != nil {
		return nil, err
	}

	divisions := toDivisionOptionResponses(projection)

	if len(divisions) > 0 {
		data, err := json.Marshal(divisions)
		if err == nil {
			_ = s.cache.Set(ctx, DivisionOptionsCacheKey, string(data), 24*time.Hour)
		}
	}

	return divisions, nil
}

func (s *Service) CreateDivision(ctx context.Context, req *DivisionCreateRequest) error {
	division := Division{
		Name: req.Name,
	}

	if err := s.repo.Create(ctx, &division); err != nil {
		if errors.Is(err, ErrDivisionAlreadyExists) {
			return apperr.AlreadyExists("division")
		}
		return err
	}

	InvalidateCache(ctx, s.cache)

	return nil
}

func (s *Service) GetDivision(ctx context.Context, id int64) (*DivisionResponse, error) {
	division, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return nil, apperr.NotFound("division")
		}
		return nil, err
	}

	result := toDivisionResponse(*division)

	return &result, nil
}

func (s *Service) UpdateDivision(ctx context.Context, id int64, req *DivisionUpdateRequest) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperr.NotFound("division")
		}
		return err
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	division := Division{
		Name:     req.Name,
		IsActive: isActive,
	}

	if err := s.repo.Update(ctx, id, &division); err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperr.NotFound("division")
		}
		if errors.Is(err, ErrDivisionAlreadyExists) {
			return apperr.AlreadyExists("division")
		}
		return err
	}

	InvalidateCache(ctx, s.cache)

	return nil
}

func (s *Service) DeleteDivision(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperr.NotFound("division")
		}
		return err
	}

	InvalidateCache(ctx, s.cache)

	return nil
}
