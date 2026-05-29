package division

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperror"
)

type DivisionService interface {
	ListDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int64, error)

	CreateDivision(ctx context.Context, req *CreateDivisionRequest) (DivisionResponse, error)

	GetDivision(ctx context.Context, id int64) (DivisionResponse, error)

	UpdateDivision(ctx context.Context, id int64, req *UpdateDivisionRequest) error

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

func (s *service) ListDivisions(ctx context.Context, params *GetDivisionParams) ([]DivisionResponse, int64, error) {

	if params == nil {

		params = &GetDivisionParams{}

	}

	params.Normalize()

	divisions, total, err := s.repo.GetAll(ctx, *params)

	if err != nil {

		return []DivisionResponse{}, 0, err

	}

	return toDivisionResponses(divisions), total, nil

}

func (s *service) CreateDivision(ctx context.Context, req *CreateDivisionRequest) (DivisionResponse, error) {

	division := Division{

		Name: req.Name,
	}

	if err := s.repo.Create(ctx, &division); err != nil {

		if errors.Is(err, ErrDivisionAlreadyExists) {

			return DivisionResponse{}, apperror.AlreadyExists("division")

		}

		return DivisionResponse{}, err

	}

	return toDivisionResponse(division), nil

}

func (s *service) GetDivision(ctx context.Context, id int64) (DivisionResponse, error) {

	division, err := s.repo.GetByID(ctx, id)

	if err != nil {

		if errors.Is(err, ErrDivisionNotFound) {

			return DivisionResponse{}, apperror.NotFound("division")

		}

		return DivisionResponse{}, err

	}

	return toDivisionResponse(*division), nil

}

func (s *service) UpdateDivision(ctx context.Context, id int64, req *UpdateDivisionRequest) error {

	existing, err := s.repo.GetByID(ctx, id)

	if err != nil {

		if errors.Is(err, ErrDivisionNotFound) {

			return apperror.NotFound("division")

		}

		return err

	}

	isActive := existing.IsActive

	if req.IsActive != nil {

		isActive = *req.IsActive

	}

	division := Division{

		Name: req.Name,

		IsActive: isActive,
	}

	if err := s.repo.Update(ctx, id, &division); err != nil {

		if errors.Is(err, ErrDivisionNotFound) {

			return apperror.NotFound("division")

		}

		if errors.Is(err, ErrDivisionAlreadyExists) {

			return apperror.AlreadyExists("division")

		}

		return err

	}

	return nil

}

func (s *service) DeleteDivision(ctx context.Context, id int64) error {

	if err := s.repo.Delete(ctx, id); err != nil {

		if errors.Is(err, ErrDivisionNotFound) {

			return apperror.NotFound("division")

		}

		return err

	}

	return nil

}
