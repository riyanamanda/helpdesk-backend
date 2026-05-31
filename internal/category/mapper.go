package category

import "github.com/riyanamanda/helpdesk-backend/internal/shared/sliceutil"

func toCategoryResponse(c Category) CategoryResponse {
	return CategoryResponse(c)
}

func toCategoryOptionResponse(c CategoryOptionProjection) CategoryOptionResponse {
	return CategoryOptionResponse(c)
}

func toCategoryResponses(categories []Category) []CategoryResponse {
	return sliceutil.Map(categories, toCategoryResponse)
}

func toCategoryOptionResponses(categories []CategoryOptionProjection) []CategoryOptionResponse {
	return sliceutil.Map(categories, toCategoryOptionResponse)
}
