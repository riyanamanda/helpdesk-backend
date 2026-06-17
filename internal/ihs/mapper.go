package ihs

func toPatientResponse(p Patient) PatientResponse {
	return PatientResponse{
		Norm:       p.Norm,
		Name:       p.Name,
		Nik:        p.Nik,
		HttpMethod: p.HttpRequest,
		GetDate:    p.GetDate,
	}
}

func toPatientResponses(patients []Patient) []PatientResponse {
	result := make([]PatientResponse, len(patients))
	for i, p := range patients {
		result[i] = toPatientResponse(p)
	}
	return result
}
