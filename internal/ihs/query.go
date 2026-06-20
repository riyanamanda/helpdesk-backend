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
		ip.getDate as get_date
	FROM ` + "`kemkes-ihs`" + `.patient ip
	JOIN master.pasien p
		ON ip.refId = p.NORM
`

func buildPatientWhere(params GetPatientParams) (string, []any) {
	var (
		where = "WHERE 1=1 AND id IS NULL AND statusRequest = 0"
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

	dir := "DESC"
	if params.SortType == "ASC" {
		dir = "ASC"
	}

	return col, dir
}
