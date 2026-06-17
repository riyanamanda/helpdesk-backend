package ihs

import "context"

type PatientService interface {
	ListPatients(ctx context.Context, params *GetPatientParams) ([]PatientResponse, int64, error)
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
