package antrian

import "fmt"

func toAntrianResponse(a Antrian) AntrianResponse {
	var waktuCheckIn *string
	if a.WaktuCheckIn.Valid {
		waktuCheckIn = &a.WaktuCheckIn.String
	}

	return AntrianResponse{
		KodeBooking:  a.KodeBooking,
		NoAntrian:    fmt.Sprintf("%s%s-%03d", a.PosAntrian, a.CaraBayar, a.No),
		Norm:         a.Norm,
		Nama:         a.Nama,
		NoKartuBpjs:  a.NoKartuBpjs,
		Dokter:       a.Dokter,
		Poli:         a.Poli,
		Status:       a.Status,
		WaktuCheckIn: waktuCheckIn,
	}
}

func toAntrianResponses(antrian []Antrian) []AntrianResponse {
	result := make([]AntrianResponse, len(antrian))
	for i, a := range antrian {
		result[i] = toAntrianResponse(a)
	}
	return result
}
