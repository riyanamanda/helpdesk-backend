package category

func toCategoryResponse(c Category) CategoryResponse {
	return CategoryResponse(c)
}

func toCategoryResponses(categories []Category) []CategoryResponse {
	result := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = toCategoryResponse(c)
	}
	return result
}
