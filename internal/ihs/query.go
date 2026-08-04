package ihs

var allowedSortColumns = map[string]string{
	"http_method": "ip.httpRequest",
	"get_date":    "ip.getDate",
}

const patientSelectBase = `
	SELECT
		ip.refId				AS norm,
		p.NAMA					AS name,
		ip.nik					AS identity_number,
		ip.httpRequest			AS http_request,
		pen_last.TANGGAL		AS last_registration,
		pen_last.nama_ruangan	AS poly
	FROM ` + "`kemkes-ihs`" + `.patient ip
	JOIN master.pasien p
		ON ip.refId = p.NORM
	LEFT JOIN (
		SELECT
			pen.NORM,
			pen.TANGGAL,
			ru.DESKRIPSI AS nama_ruangan,
			ROW_NUMBER() OVER (
				PARTITION BY pen.NORM
				ORDER BY pen.TANGGAL DESC, pen.NOMOR DESC
			) AS rn
		FROM pendaftaran.pendaftaran pen
		JOIN pendaftaran.kunjungan kun
			ON kun.NOPEN = pen.NOMOR
		LEFT JOIN pendaftaran.tujuan_pasien tu
			ON tu.NOPEN = kun.NOPEN
		LEFT JOIN master.ruangan ru
			ON ru.ID = tu.RUANGAN
		WHERE pen.STATUS IN (1,2)
		AND kun.STATUS IN (1,2)
	) pen_last
		ON pen_last.NORM = ip.refId
	AND pen_last.rn = 1
`

func buildPatientWhere(params GetPatientParams) (string, []any) {
	var (
		where = "WHERE 1=1 AND ip.id IS NULL AND ip.statusRequest = 0"
		args  []any
	)

	if params.Search != "" {
		like := "%" + params.Search + "%"
		args = append(args, like, like)
		where += " AND (ip.refId LIKE ? OR ip.nik LIKE ?)"
	}

	if params.HttpMethod != "" {
		args = append(args, params.HttpMethod)
		where += " AND ip.httpRequest = ?"
	}

	return where, args
}

func buildPatientSort(params GetPatientParams) (string, string) {
	col, ok := allowedSortColumns[params.SortBy]
	if !ok {
		col = "ip.getDate"
	}

	dir := "ASC"
	if params.SortType == "DESC" {
		dir = "DESC"
	}

	return col, dir
}
