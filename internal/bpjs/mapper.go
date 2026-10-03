package bpjs

func toPesertaResponse(vclaim VclaimResponse) *PesertaResponse {
	return &PesertaResponse{
		Nama:     vclaim.Peserta.Nama,
		NIK:      vclaim.Peserta.NIK,
		NoKartu:  vclaim.Peserta.NoKartu,
		Sex:      vclaim.Peserta.Sex,
		TglLahir: vclaim.Peserta.TglLahir,
		TglTAT:   vclaim.Peserta.TglTAT,
		TglTMT:   vclaim.Peserta.TglTMT,
	}
}
