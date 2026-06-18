package ihs

import "time"

type Patient struct {
	Norm           string    `db:"norm"`
	Name           string    `db:"name"`
	IdentityNumber string    `db:"identity_number"`
	HttpRequest    string    `db:"http_request"`
	GetDate        time.Time `db:"get_date"`
}

type PatientDetail struct {
	Norm          string       `db:"norm"`
	Name          string       `db:"name"`
	BirthPlace    string       `db:"birt_place"`
	BirthDate     string       `db:"birt_date"`
	MaritalStatus string       `db:"marital_status"`
	Citizenship   string       `db:"citizenship"`
	Status        bool         `db:"status"`
	IdentityCard  IdentityCard `db:"identity_card"`
}

type IdentityCard struct {
	IdentityNumber string `db:"identity_number"`
	Address        string `db:"address"`
	RT             string `db:"rt"`
	RW             string `db:"rw"`
	Province       string `db:"province"`
	City           string `db:"city"`
	District       string `db:"district"`
	SubDistrict    string `db:"sub_district"`
}
