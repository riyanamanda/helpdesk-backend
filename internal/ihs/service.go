package ihs

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

type repository interface {
	GetPatients(ctx context.Context, params GetPatientParams) ([]PatientProjection, int64, error)
	GetPatientDetail(ctx context.Context, NORM string) (*PatientDetailProjection, error)
	UpdatePatientMethod(ctx context.Context, NORM string) error
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListPatients(ctx context.Context, params *GetPatientParams) ([]PatientResponse, int64, error) {
	if params == nil {
		params = &GetPatientParams{}
	}
	params.Normalize()

	patients, total, err := s.repo.GetPatients(ctx, *params)
	if err != nil {
		return nil, 0, err
	}

	return toPatientResponses(patients), total, nil
}

func (s *Service) GetPatientByNORM(ctx context.Context, NORM string) (*PatientDetailResponse, error) {
	patient, err := s.repo.GetPatientDetail(ctx, NORM)
	if err != nil {
		if errors.Is(err, ErrPatientNotFound) {
			return nil, apperr.NotFound("patient")
		}
		return nil, err
	}

	result := toPatientDetailResponse(*patient)
	return &result, nil
}

func (s *Service) UpdatePatientMethodByNORM(ctx context.Context, NORM string) error {
	if err := s.repo.UpdatePatientMethod(ctx, NORM); err != nil {
		if errors.Is(err, ErrPatientNotFound) {
			return apperr.NotFound("patient")
		}
		if errors.Is(err, ErrPatientNotEligible) {
			return apperr.Conflict("patient has already been submitted or is currently being processed")
		}
		return err
	}

	return nil
}
