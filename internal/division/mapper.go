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
	result := make([]DivisionResponse, 0, len(divisions))
	for _, d := range divisions {
		result = append(result, toDivisionResponse(d))
	}
	return result
}
