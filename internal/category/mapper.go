package category

import "github.com/riyanamanda/helpdesk-backend/internal/shared/collection"

func toCategoryResponse(c Category) CategoryResponse {
	return CategoryResponse(c)
}

func toCategoryResponses(categories []Category) []CategoryResponse {
	return collection.MapSlice(categories, toCategoryResponse)
}
