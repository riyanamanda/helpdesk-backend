package division

import "github.com/riyanamanda/helpdesk-backend/internal/shared/collection"

func toDivisionResponse(d Division) DivisionResponse {
	return DivisionResponse(d)
}

func toDivisionResponses(divisions []Division) []DivisionResponse {
	return collection.MapSlice(divisions, toDivisionResponse)
}
