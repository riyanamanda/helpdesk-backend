package ihs

type PatientDetailProjection struct {
	Norm           string `db:"norm"`
	Name           string `db:"name"`
	BirthPlace     string `db:"birth_place"`
	BirthDate      string `db:"birth_date"`
	Gender         string `db:"gender"`
	MaritalStatus  string `db:"marital_status"`
	Citizenship    string `db:"citizenship"`
	Status         bool   `db:"status"`
	IdentityNumber string `db:"identity_number"`
	Address        string `db:"address"`
	RT             string `db:"rt"`
	RW             string `db:"rw"`
	Provnice       string `db:"province"`
	City           string `db:"city"`
	District       string `db:"district"`
	SubDistrict    string `db:"sub_district"`
}
