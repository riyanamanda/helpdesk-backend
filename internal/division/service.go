package division

import (
	"context"
	"errors"

	apperrors "github.com/riyanamanda/helpdesk-backend/internal/shared/errors"
)

type DivisionService interface {
	FetchAllDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int64, error)
	RegisterDivision(ctx context.Context, req *CreateDivisionRequest) (DivisionResponse, error)
	FindDivisionByID(ctx context.Context, id int64) (DivisionResponse, error)
	EditDivision(ctx context.Context, id int64, req *UpdateDivisionRequest) error
	DeleteDivision(ctx context.Context, id int64) error
}

type service struct {
	repo DivisionRepository
}

func NewDivisionService(repo DivisionRepository) DivisionService {
	return &service{
		repo: repo,
	}
}

func (svc *service) FetchAllDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int64, error) {
	if params == nil {
		params = &GetDivisionParams{}
	}

	params.Normalize()

	divisions, total, err := svc.repo.GetAll(ctx, *params)
	if err != nil {
		return []DivisionResponse{}, 0, err
	}

	return toDivisionResponses(divisions), total, nil
}

func (svc *service) RegisterDivision(ctx context.Context, req *CreateDivisionRequest) (DivisionResponse, error) {
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

func (svc *service) FindDivisionByID(ctx context.Context, id int64) (DivisionResponse, error) {
	division, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return DivisionResponse{}, apperrors.NotFound("division")
		}
		return DivisionResponse{}, err
	}

	return toDivisionResponse(*division), nil
}

func (svc *service) EditDivision(ctx context.Context, id int64, req *UpdateDivisionRequest) error {
	existing, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperrors.NotFound("division")
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
	if err := svc.repo.Update(ctx, id, &division); err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperrors.NotFound("division")
		}
		if errors.Is(err, ErrDivisionAlreadyExists) {
			return apperrors.AlreadyExists("division")
		}
		return err
	}

	return nil
}

func (svc *service) DeleteDivision(ctx context.Context, id int64) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrDivisionNotFound) {
			return apperrors.NotFound("division")
		}
		return err
	}

	return nil
}
