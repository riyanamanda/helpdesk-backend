package ihs

import "fmt"

var allowedSortColumns = map[string]string{
	"http_method":       "ip.httpRequest",
	"get_date":          "ip.getDate",
	"last_registration": "pen_last.TANGGAL",
}

const patientFromBase = `
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

const patientSelectBase = `
	SELECT
		ip.refId				AS norm,
		p.NAMA					AS name,
		ip.nik					AS identity_number,
		ip.httpRequest			AS http_request,
		pen_last.TANGGAL		AS last_registration,
		pen_last.nama_ruangan	AS poly
` + patientFromBase

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

	if params.StartDate != "" {
		args = append(args, params.StartDate)
		where += " AND pen_last.TANGGAL >= ?"
	}

	if params.EndDate != "" {
		args = append(args, params.EndDate)
		where += " AND pen_last.TANGGAL <= ?"
	}

	return where, args
}

func buildPatientSort(params GetPatientParams) string {
	col, ok := allowedSortColumns[params.SortBy]
	if !ok {
		col = "pen_last.TANGGAL"
	}

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	return fmt.Sprintf("%s %s, (ip.httpRequest = 'GET') DESC", col, dir)
}
