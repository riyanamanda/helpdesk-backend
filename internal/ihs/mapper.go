package ihs

func toPatientResponse(p Patient) PatientResponse {
	return PatientResponse{
		Norm:           p.Norm,
		Name:           p.Name,
		IdentityNumber: p.IdentityNumber,
		HttpMethod:     p.HttpRequest,
		GetDate:        p.GetDate,
	}
}

func toPatientDetailResponse(p PatientDetailProjection) PatientDetailResponse {
	return PatientDetailResponse{
		Norm:          p.Norm,
		Name:          p.Name,
		BirthPlace:    p.BirthPlace,
		BirthDate:     p.BirthDate,
		MaritalStatus: p.MaritalStatus,
		Citizenship:   p.Citizenship,
		Status:        p.Status,
		IdentityCard: IdentityCardResponse{
			IdentityNumber: p.IdentityNumber,
			Address:        p.Address,
			RT:             p.RT,
			RW:             p.RW,
			Province:       p.Provnice,
			City:           p.City,
			District:       p.District,
			SubDistrict:    p.SubDistrict,
		},
	}
}

func toPatientResponses(patients []Patient) []PatientResponse {
	result := make([]PatientResponse, len(patients))
	for i, p := range patients {
		result[i] = toPatientResponse(p)
	}
	return result
}
