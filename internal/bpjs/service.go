package bpjs

import (
	"context"
	"encoding/json"
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetPesertaByNIK(ctx context.Context, nik string) (*PesertaResponse, error) {
	plainJson, err := s.client.GetPeserta(ctx, nik)
	if err != nil {
		return nil, err
	}

	var patient VclaimResponse
	if err := json.Unmarshal(plainJson, &patient); err != nil {
		return nil, err
	}

	return toPesertaResponse(patient), nil
}
