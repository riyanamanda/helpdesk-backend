package antrian

import "github.com/riyanamanda/helpdesk-backend/internal/shared/pagination"

type GetAntrianParams struct {
	pagination.Params
	Norm string `query:"norm"`
}

func (p *GetAntrianParams) Normalize() {
	p.Params.Normalize()
}

type AntrianResponse struct {
	KodeBooking  int64   `json:"kode_booking"`
	NoAntrian    string  `json:"no_antrian"`
	Norm         string  `json:"norm"`
	Nama         string  `json:"nama"`
	NoKartuBpjs  string  `json:"no_kartu_bpjs"`
	Dokter       string  `json:"dokter"`
	Poli         string  `json:"poli"`
	Status       int     `json:"status"`
	WaktuCheckIn *string `json:"waktu_check_in"`
}
