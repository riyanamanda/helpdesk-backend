package antrian

import (
	"context"

	"github.com/riyanamanda/helpdesk-backend/internal/simgos"
)

type AntrianService interface {
	ListAntrian(ctx context.Context, params *GetAntrianParams) ([]AntrianResponse, int64, error)
	CheckIn(ctx context.Context, kodeBooking int64) error
}

type service struct {
	repo   AntrianRepository
	simgos *simgos.AntrolClient
}

func NewAntrianService(repo AntrianRepository, simgos *simgos.AntrolClient) AntrianService {
	return &service{
		repo:   repo,
		simgos: simgos,
	}
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
	return s.simgos.CheckIn(ctx, kodeBooking)
}
