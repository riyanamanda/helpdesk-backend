package category

func toCategoryResponse(c Category) CategoryResponse {
	return CategoryResponse(c)
}

func toCategoryResponses(categories []Category) []CategoryResponse {
	result := make([]CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryResponse(c))
	}
	return result
}
