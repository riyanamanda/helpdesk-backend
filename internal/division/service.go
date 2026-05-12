package division

import (
	"context"
	"errors"
	"log/slog"

	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type DivisionService interface {
	GetDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int, error)
	Create(ctx context.Context, req *CreateDivisionRequest) (DivisionResponse, error)
	GetByID(ctx context.Context, id int64) (DivisionResponse, error)
	Update(ctx context.Context, id int64, req *UpdateDivisionRequest) (DivisionResponse, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo DivisionRepository
}

func NewDivisionService(repo DivisionRepository) DivisionService {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int, error) {
	if params == nil {
		params = &GetDivisionParams{}
	}

	page, limit, _ := params.Normalize()
	params.Page = page
	params.Limit = limit

	divisions, total, err := svc.repo.List(ctx, *params)
	if err != nil {
		slog.Error("list division failed", "error", err)
		return []DivisionResponse{}, 0, err
	}

	return toDivisionResponses(divisions), total, nil
}

func (svc *service) Create(ctx context.Context, req *CreateDivisionRequest) (DivisionResponse, error) {
	division := Division{
		Name: req.Name,
	}

	if err := svc.repo.Create(ctx, &division); err != nil {
		if errors.Is(err, ErrDivisionAlreadyExists) {
			return DivisionResponse{}, apperrors.AlreadyExists("division")
		}
		return DivisionResponse{}, err
	}

	return toDivisionResponse(division), nil
}

func (svc *service) GetByID(ctx context.Context, id int64) (DivisionResponse, error) {
	division, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return DivisionResponse{}, apperrors.NotFound("division")
		}
		return DivisionResponse{}, err
	}

	return toDivisionResponse(*division), nil
}

func (svc *service) Update(ctx context.Context, id int64, req *UpdateDivisionRequest) (DivisionResponse, error) {
	division := Division{Name: req.Name}
	if err := svc.repo.Update(ctx, id, &division); err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return DivisionResponse{}, apperrors.NotFound("division")
		}
		if errors.Is(err, ErrDivisionAlreadyExists) {
			return DivisionResponse{}, apperrors.AlreadyExists("division")
		}
		return DivisionResponse{}, err
	}

	return toDivisionResponse(division), nil
}

func (svc *service) Delete(ctx context.Context, id int64) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperrors.NotFound("division")
		}
		return err
	}

	return nil
}
