package ihs

import (
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type PatientResponse struct {
	Norm       string    `json:"norm"`
	Name       string    `json:"name"`
	Nik        string    `json:"nik"`
	HttpMethod string    `json:"http_method"`
	GetDate    time.Time `json:"get_date"`
}

type GetPatientParams struct {
	pagination.Params
	HttpMethod string `query:"http_method"`
}

func (p *GetPatientParams) Normalize() {
	p.Params.Normalize()
}
