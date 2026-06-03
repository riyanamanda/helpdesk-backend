package division

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
)

type DivisionService interface {
	ListDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int64, error)
	ListOptions(ctx context.Context) ([]DivisionOptionResponse, error)
	CreateDivision(ctx context.Context, req *DivisionCreateRequest) (*DivisionResponse, error)
	GetDivision(ctx context.Context, id int64) (*DivisionResponse, error)
	UpdateDivision(ctx context.Context, id int64, req *DivisionUpdateRequest) error
	DeleteDivision(ctx context.Context, id int64) error
}

type service struct {
	repo  DivisionRepository
	cache cache.Cache
}

func NewDivisionService(repo DivisionRepository, cache cache.Cache) DivisionService {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) ListDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int64, error) {
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

func (s *service) ListOptions(ctx context.Context) ([]DivisionOptionResponse, error) {
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

func (s *service) CreateDivision(ctx context.Context, req *DivisionCreateRequest) (*DivisionResponse, error) {
	division := Division{
		Name: req.Name,
	}

	if err := s.repo.Create(ctx, &division); err != nil {
		if errors.Is(err, ErrDivisionAlreadyExists) {
			return nil, apperr.AlreadyExists("division")
		}
		return nil, err
	}

	InvalidateCache(ctx, s.cache)

	result := toDivisionResponse(division)

	return &result, nil
}

func (s *service) GetDivision(ctx context.Context, id int64) (*DivisionResponse, error) {
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

func (s *service) UpdateDivision(ctx context.Context, id int64, req *DivisionUpdateRequest) error {
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

func (s *service) DeleteDivision(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperr.NotFound("division")
		}
		return err
	}

	InvalidateCache(ctx, s.cache)

	return nil
}
