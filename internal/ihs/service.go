package ihs

import (
	"context"
	"errors"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/apperr"
)

type PatientService interface {
	ListPatients(ctx context.Context, params *GetPatientParams) ([]PatientResponse, int64, error)
	GetPatientByNORM(ctx context.Context, NORM string) (*PatientDetailResponse, error)
	UpdatePatientMethodByNORM(ctx context.Context, NORM string) error
}

type service struct {
	repo PatientRepository
}

func NewPatientService(repo PatientRepository) PatientService {
	return &service{
		repo: repo,
	}
}

func (s *service) ListPatients(ctx context.Context, params *GetPatientParams) ([]PatientResponse, int64, error) {
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

func (s *service) GetPatientByNORM(ctx context.Context, NORM string) (*PatientDetailResponse, error) {
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

func (s *service) UpdatePatientMethodByNORM(ctx context.Context, NORM string) error {
	if err := s.repo.UpdatePatientMethod(ctx, NORM); err != nil {
		if errors.Is(err, ErrPatientNotFound) {
			return apperr.NotFound("patient")
		}
		return err
	}

	return nil
}
