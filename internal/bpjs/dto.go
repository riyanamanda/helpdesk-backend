package bpjs

type BPJSResponse struct {
	MetaData struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"metaData"`
	Response string `json:"response"`
}

type PesertaResponse struct {
	Nama     string `json:"nama"`
	NIK      string `json:"nik"`
	NoKartu  string `json:"no_kartu"`
	Sex      string `json:"sex"`
	TglLahir string `json:"tgl_lahir"`
	TglTAT   string `json:"tgl_tat"`
	TglTMT   string `json:"tgl_tmt"`
}

type VclaimResponse struct {
	Peserta struct {
		Nama     string `json:"nama"`
		NIK      string `json:"nik"`
		NoKartu  string `json:"noKartu"`
		Sex      string `json:"sex"`
		TglLahir string `json:"tglLahir"`
		TglTAT   string `json:"tglTAT"`
		TglTMT   string `json:"tglTMT"`
	} `json:"peserta"`
}
