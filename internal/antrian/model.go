package antrian

import "database/sql"

type Antrian struct {
	KodeBooking  int64          `db:"kode_booking"`
	PosAntrian   string         `db:"pos_antrian"`
	CaraBayar    string         `db:"cara_bayar"`
	No           int            `db:"no"`
	Norm         string         `db:"norm"`
	Nama         string         `db:"nama"`
	NoKartuBpjs  string         `db:"no_kartu_bpjs"`
	Dokter       string         `db:"dokter"`
	Poli         string         `db:"poli"`
	Status       int            `db:"status"`
	WaktuCheckIn sql.NullString `db:"waktu_check_in"`
}
