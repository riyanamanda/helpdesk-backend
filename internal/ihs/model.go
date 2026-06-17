package ihs

import "time"

type Patient struct {
	Norm        string    `db:"norm"`
	Name        string    `db:"name"`
	Nik         string    `db:"nik"`
	HttpRequest string    `db:"http_request"`
	GetDate     time.Time `db:"get_date"`
}
