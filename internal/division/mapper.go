package division

func toDivisionResponse(d Division) DivisionResponse {
	return DivisionResponse(d)
}

func toDivisionOptionResponse(d DivisionOptionProjection) DivisionOptionResponse {
	return DivisionOptionResponse(d)
}

func toDivisionResponses(divisions []Division) []DivisionResponse {
	result := make([]DivisionResponse, len(divisions))
	for i, d := range divisions {
		result[i] = toDivisionResponse(d)
	}
	return result
}

func toDivisionOptionResponses(divisions []DivisionOptionProjection) []DivisionOptionResponse {
	result := make([]DivisionOptionResponse, len(divisions))
	for i, d := range divisions {
		result[i] = toDivisionOptionResponse(d)
	}
	return result
}
