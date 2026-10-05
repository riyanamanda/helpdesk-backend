package bpjs

// BPJSResponse adalah wrapper awal untuk memegang string terenkripsi dari VClaim
type BPJSResponse struct {
	MetaData struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"metaData"`
	Response string `json:"response"`
}

// -----------------------------------------------------------------------------
// 1. DTO Publik (dikirim ke Client Application dengan format snake_case)
// -----------------------------------------------------------------------------

type PesertaResponse struct {
	NoKartu       string        `json:"no_kartu"`
	NIK           string        `json:"nik"`
	Nama          string        `json:"nama"`
	Pisa          string        `json:"pisa"`
	Sex           string        `json:"sex"`
	TglLahir      string        `json:"tgl_lahir"`
	TglCetakKartu string        `json:"tgl_cetak_kartu"`
	TglTAT        string        `json:"tgl_tat"`
	TglTMT        string        `json:"tgl_tmt"`
	MR            MedicalRecord `json:"mr"`
	StatusPeserta CodeDesc      `json:"status_peserta"`
	ProvUmum      Provider      `json:"prov_umum"`
	JenisPeserta  CodeDesc      `json:"jenis_peserta"`
	HakKelas      CodeDesc      `json:"hak_kelas"`
	Umur          Umur          `json:"umur"`
	Informasi     Informasi     `json:"informasi"`
	COB           COB           `json:"cob"`
}

type MedicalRecord struct {
	NoMR      *string `json:"no_mr"`
	NoTelepon *string `json:"no_telepon"`
}

type CodeDesc struct {
	Kode       string `json:"kode"`
	Keterangan string `json:"keterangan"`
}

type Provider struct {
	KdProvider string `json:"kd_provider"`
	NmProvider string `json:"nm_provider"`
}

type Umur struct {
	UmurSekarang      string `json:"umur_sekarang"`
	UmurSaatPelayanan string `json:"umur_saat_pelayanan"`
}

type Informasi struct {
	Dinsos      *string `json:"dinsos"`
	ProlanisPRB *string `json:"prolanis_prb"`
	NoSKTM      *string `json:"no_sktm"`
	ESEP        *string `json:"e_sep"`
}

type COB struct {
	NoAsuransi *string `json:"no_asuransi"`
	NmAsuransi *string `json:"nm_asuransi"`
	TglTMT     *string `json:"tgl_tmt"`
	TglTAT     *string `json:"tgl_tat"`
}

// -----------------------------------------------------------------------------
// 2. Struct Internal (untuk Unmarshal JSON Mentah / Decrypted BPJS - camelCase)
// -----------------------------------------------------------------------------

type VclaimResponse struct {
	Peserta struct {
		NoKartu       string `json:"noKartu"`
		NIK           string `json:"nik"`
		Nama          string `json:"nama"`
		Pisa          string `json:"pisa"`
		Sex           string `json:"sex"`
		TglLahir      string `json:"tglLahir"`
		TglCetakKartu string `json:"tglCetakKartu"`
		TglTAT        string `json:"tglTAT"`
		TglTMT        string `json:"tglTMT"`
		MR            struct {
			NoMR      *string `json:"noMR"`
			NoTelepon *string `json:"noTelepon"`
		} `json:"mr"`
		StatusPeserta struct {
			Kode       string `json:"kode"`
			Keterangan string `json:"keterangan"`
		} `json:"statusPeserta"`
		ProvUmum struct {
			KdProvider string `json:"kdProvider"`
			NmProvider string `json:"nmProvider"`
		} `json:"provUmum"`
		JenisPeserta struct {
			Kode       string `json:"kode"`
			Keterangan string `json:"keterangan"`
		} `json:"jenisPeserta"`
		HakKelas struct {
			Kode       string `json:"kode"`
			Keterangan string `json:"keterangan"`
		} `json:"hakKelas"`
		Umur struct {
			UmurSekarang      string `json:"umurSekarang"`
			UmurSaatPelayanan string `json:"umurSaatPelayanan"`
		} `json:"umur"`
		Informasi struct {
			Dinsos      *string `json:"dinsos"`
			ProlanisPRB *string `json:"prolanisPRB"`
			NoSKTM      *string `json:"noSKTM"`
			ESEP        *string `json:"eSEP"`
		} `json:"informasi"`
		COB struct {
			NoAsuransi *string `json:"noAsuransi"`
			NmAsuransi *string `json:"nmAsuransi"`
			TglTMT     *string `json:"tglTMT"`
			TglTAT     *string `json:"tglTAT"`
		} `json:"cob"`
	} `json:"peserta"`
}
