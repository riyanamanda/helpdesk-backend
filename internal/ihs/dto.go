package ihs

import (
	"time"

	"github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"
)

type PatientResponse struct {
	Norm           string    `json:"norm"`
	Name           string    `json:"name"`
	IdentityNumber string    `json:"identity_number"`
	HttpMethod     string    `json:"http_method"`
	GetDate        time.Time `json:"get_date"`
}

type PatientDetailResponse struct {
	Norm          string               `json:"norm"`
	Name          string               `json:"name"`
	BirthPlace    *string              `json:"birth_place"`
	BirthDate     string               `json:"birth_date"`
	MaritalStatus string               `json:"marital_status"`
	Citizenship   string               `json:"citizenship"`
	Status        bool                 `json:"status"`
	IdentityCard  IdentityCardResponse `json:"identity_card"`
}

type IdentityCardResponse struct {
	IdentityNumber string `json:"identity_number"`
	Address        string `json:"address"`
	RT             string `json:"rt"`
	RW             string `json:"rw"`
	Province       string `json:"province"`
	City           string `json:"city"`
	District       string `json:"district"`
	SubDistrict    string `json:"sub_district"`
}

type GetPatientParams struct {
	pagination.Params
	HttpMethod string `query:"http_method"`
}

func (p *GetPatientParams) Normalize() {
	p.Params.Normalize()
}
