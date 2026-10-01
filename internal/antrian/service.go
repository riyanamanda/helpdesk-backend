package antrian

import (
	"context"

	"github.com/riyanamanda/helpdesk-backend/internal/simgos"
)

type repository interface {
	GetAntrian(ctx context.Context, params GetAntrianParams) ([]Antrian, int64, error)
}

type Service struct {
	repo   repository
	simgos *simgos.AntrolClient
}

func NewService(repo repository, simgos *simgos.AntrolClient) *Service {
	return &Service{
		repo:   repo,
		simgos: simgos,
	}
}

func (s *Service) ListAntrian(ctx context.Context, params *GetAntrianParams) ([]AntrianResponse, int64, error) {
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

func (s *Service) CheckIn(ctx context.Context, kodeBooking int64) error {
	return s.simgos.CheckIn(ctx, kodeBooking)
}
