package ticket

import (
	"context"
)

type TicketService interface {
	FetchAllTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int, error)
}

type service struct {
	repo TicketRepository
}

func NewTicketService(repo TicketRepository) TicketService {
	return &service{
		repo: repo,
	}
}

func (s *service) FetchAllTickets(ctx context.Context, params *GetTicketParams) ([]TicketResponse, int, error) {
	if params == nil {
		params = &GetTicketParams{}
	}

	page, limit, _ := params.Normalize()
	params.Page = page
	params.Limit = limit

	tickets, total, err := s.repo.GetAll(ctx, *params)
	if err != nil {
		return []TicketResponse{}, 0, err
	}

	return toTicketResponses(tickets), total, nil
}
