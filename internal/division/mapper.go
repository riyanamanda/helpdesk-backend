package division

func toDivisionResponse(d Division) DivisionResponse {
	return DivisionResponse{
		ID:        d.ID,
		Name:      d.Name,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func toDivisionResponses(divisions []Division) []DivisionResponse {
	result := make([]DivisionResponse, len(divisions))
	for i, d := range divisions {
		result[i] = toDivisionResponse(d)
	}
	return result
}
