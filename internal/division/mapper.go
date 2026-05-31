package division

import "github.com/riyanamanda/helpdesk-backend/internal/shared/sliceutil"

func toDivisionResponse(d Division) DivisionResponse {
	return DivisionResponse(d)
}

func toDivisionOptionResponse(d DivisionOptionProjection) DivisionOptionResponse {
	return DivisionOptionResponse(d)
}

func toDivisionResponses(divisions []Division) []DivisionResponse {
	return sliceutil.Map(divisions, toDivisionResponse)
}

func toDivisionOptionResponses(divisions []DivisionOptionProjection) []DivisionOptionResponse {
	return sliceutil.Map(divisions, toDivisionOptionResponse)
}
