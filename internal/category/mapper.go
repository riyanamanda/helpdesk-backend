package category

func toCategoryResponse(c Category) CategoryResponse {
	return CategoryResponse(c)
}

func toCategoryOptionResponse(c CategoryOptionProjection) CategoryOptionResponse {
	return CategoryOptionResponse(c)
}

func toCategoryResponses(categories []Category) []CategoryResponse {
	result := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = toCategoryResponse(c)
	}
	return result
}

func toCategoryOptionResponses(categories []CategoryOptionProjection) []CategoryOptionResponse {
	result := make([]CategoryOptionResponse, len(categories))
	for i, c := range categories {
		result[i] = toCategoryOptionResponse(c)
	}
	return result
}
