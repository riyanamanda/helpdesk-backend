package bpjs

func toPesertaResponse(vclaim VclaimResponse) *PesertaResponse {
	p := vclaim.Peserta

	return &PesertaResponse{
		NoKartu:       p.NoKartu,
		NIK:           p.NIK,
		Nama:          p.Nama,
		Pisa:          p.Pisa,
		Sex:           p.Sex,
		TglLahir:      p.TglLahir,
		TglCetakKartu: p.TglCetakKartu,
		TglTAT:        p.TglTAT,
		TglTMT:        p.TglTMT,
		MR: MedicalRecord{
			NoMR:      p.MR.NoMR,
			NoTelepon: p.MR.NoTelepon,
		},
		StatusPeserta: CodeDesc{
			Kode:       p.StatusPeserta.Kode,
			Keterangan: p.StatusPeserta.Keterangan,
		},
		ProvUmum: Provider{
			KdProvider: p.ProvUmum.KdProvider,
			NmProvider: p.ProvUmum.NmProvider,
		},
		JenisPeserta: CodeDesc{
			Kode:       p.JenisPeserta.Kode,
			Keterangan: p.JenisPeserta.Keterangan,
		},
		HakKelas: CodeDesc{
			Kode:       p.HakKelas.Kode,
			Keterangan: p.HakKelas.Keterangan,
		},
		Umur: Umur{
			UmurSekarang:      p.Umur.UmurSekarang,
			UmurSaatPelayanan: p.Umur.UmurSaatPelayanan,
		},
		Informasi: Informasi{
			Dinsos:      p.Informasi.Dinsos,
			ProlanisPRB: p.Informasi.ProlanisPRB,
			NoSKTM:      p.Informasi.NoSKTM,
			ESEP:        p.Informasi.ESEP,
		},
		COB: COB{
			NoAsuransi: p.COB.NoAsuransi,
			NmAsuransi: p.COB.NmAsuransi,
			TglTMT:     p.COB.TglTMT,
			TglTAT:     p.COB.TglTAT,
		},
	}
}
