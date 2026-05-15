package division

func toDivisionResponse(d Division) DivisionResponse {
	return DivisionResponse(d)
}

func toDivisionResponses(divisions []Division) []DivisionResponse {
	result := make([]DivisionResponse, len(divisions))
	for i, d := range divisions {
		result[i] = toDivisionResponse(d)
	}
	return result
}
