package antrian

import "context"

type AntrianService interface {
	ListAntrian(ctx context.Context, params *GetAntrianParams) ([]AntrianResponse, int64, error)
	CheckIn(ctx context.Context, kodeBooking int64) error
}

type service struct {
	repo   AntrianRepository
	antrol *antrolClient
}

func NewAntrianService(repo AntrianRepository, antrol *antrolClient) AntrianService {
	return &service{repo: repo, antrol: antrol}
}

func (s *service) ListAntrian(ctx context.Context, params *GetAntrianParams) ([]AntrianResponse, int64, error) {
	if params == nil {
		params = &GetAntrianParams{}
	}
	params.Normalize()

	antrian, total, err := s.repo.GetAntrian(ctx, *params)
	if err != nil {
		return nil, 0, err
	}

	return toAntrianResponses(antrian), total, nil
}

func (s *service) CheckIn(ctx context.Context, kodeBooking int64) error {
	return s.antrol.checkIn(ctx, kodeBooking)
}
