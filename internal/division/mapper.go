package division

import "github.com/riyanamanda/helpdesk-backend/internal/shared/sliceutil"

func toDivisionResponse(d Division) DivisionResponse {
	return DivisionResponse(d)
}

func toDivisionResponses(divisions []Division) []DivisionResponse {
	return sliceutil.Map(divisions, toDivisionResponse)
}
