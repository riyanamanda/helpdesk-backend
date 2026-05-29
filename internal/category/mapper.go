package category

import "github.com/riyanamanda/helpdesk-backend/internal/shared/sliceutil"

func toCategoryResponse(c Category) CategoryResponse {
	return CategoryResponse(c)
}

func toCategoryResponses(categories []Category) []CategoryResponse {
	return sliceutil.Map(categories, toCategoryResponse)
}
