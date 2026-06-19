package antrian

const antrianSelectBase = `
	SELECT
		r.ID             AS kode_booking,
		r.POS_ANTRIAN    AS pos_antrian,
		r.CARABAYAR      AS cara_bayar,
		r.NO             AS no,
		r.NORM           AS norm,
		r.NAMA           AS nama,
		r.NO_KARTU_BPJS  AS no_kartu_bpjs,
		COALESCE(d.NAMA, '') AS dokter,
		r.POLI_BPJS      AS poli,
		r.STATUS         AS status,
		r.WAKTU_CHECK_IN AS waktu_check_in
	FROM regonline.reservasi r
	LEFT JOIN regonline.dokter d ON d.KODE = r.DOKTER
`

func buildAntrianWhere(params GetAntrianParams) (string, []any) {
	where := "WHERE r.TANGGALKUNJUNGAN = CURDATE() AND r.JENIS_APLIKASI = 2"
	var args []any

	if params.Norm != "" {
		where += " AND r.STATUS IN (0, 1, 2, 99) AND r.NORM = ?"
		args = append(args, params.Norm)
	} else {
		where += " AND r.STATUS = 1"
	}

	return where, args
}
