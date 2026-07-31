package ihs

var allowedSortColumns = map[string]string{
	"http_method": "ip.httpRequest",
	"get_date":    "ip.getDate",
}

const patientSelectBase = `
    SELECT
        ip.refId as norm,
        p.NAMA as name,
        ip.nik as identity_number,
        ip.httpRequest as http_request,
        ip.getDate as get_date,
        (
            SELECT pen.TANGGAL
            FROM pendaftaran.pendaftaran pen
            JOIN pendaftaran.kunjungan kun
                ON kun.NOPEN = pen.NOMOR
            WHERE pen.NORM = ip.refId
              AND pen.STATUS = 2
              AND kun.STATUS = 2
            ORDER BY pen.TANGGAL DESC, pen.NOMOR DESC
            LIMIT 1
        ) as last_registration
    FROM ` + "`kemkes-ihs`" + `.patient ip
    JOIN master.pasien p
        ON ip.refId = p.NORM
`

func buildPatientWhere(params GetPatientParams) (string, []any) {
	var (
		where = "WHERE 1=1 AND ip.id IS NULL AND ip.statusRequest = 0"
		args  []any
	)

	if params.StartDate != "" && params.EndDate != "" {
		where += ` AND (
			SELECT pen.TANGGAL
			FROM pendaftaran.pendaftaran pen
			JOIN pendaftaran.kunjungan kun
				ON kun.NOPEN = pen.NOMOR
			WHERE pen.NORM = ip.refId
			  AND pen.STATUS = 2
			  AND kun.STATUS = 2
			ORDER BY pen.TANGGAL DESC, pen.NOMOR DESC
			LIMIT 1
		) BETWEEN ? AND ?`

		args = append(args, params.StartDate, params.EndDate)
	}

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
