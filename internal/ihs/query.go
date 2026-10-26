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
            MAX(pen.TANGGAL) AS TANGGAL,
            MAX(ru.DESKRIPSI) AS nama_ruangan
        FROM pendaftaran.pendaftaran pen
        JOIN (
            SELECT pmax.NORM, MAX(pmax.NOMOR) AS max_nomor
            FROM pendaftaran.pendaftaran pmax
            WHERE pmax.STATUS IN (1, 2)
              AND pmax.NORM IN (
                  SELECT ip_sub.refId 
                  FROM ` + "`kemkes-ihs`" + `.patient ip_sub 
                  WHERE ip_sub.id IS NULL AND ip_sub.statusRequest = 0
              )
            GROUP BY pmax.NORM
        ) latest
            ON pen.NORM = latest.NORM 
            AND pen.NOMOR = latest.max_nomor
        JOIN pendaftaran.kunjungan kun
            ON kun.NOPEN = pen.NOMOR
        LEFT JOIN pendaftaran.tujuan_pasien tu
            ON tu.NOPEN = kun.NOPEN
        LEFT JOIN master.ruangan ru
            ON ru.ID = tu.RUANGAN
        WHERE kun.STATUS IN (1, 2)
        GROUP BY pen.NORM
    ) pen_last
        ON pen_last.NORM = ip.refId
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

	httpMethod := params.HttpMethod

	if params.Search != "" {
		if httpMethod == "" || httpMethod == "GET" {
			httpMethod = "ALL"
		}
	} else {
		if httpMethod == "" {
			httpMethod = "GET"
		}
	}

	if httpMethod != "ALL" {
		args = append(args, httpMethod)
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
