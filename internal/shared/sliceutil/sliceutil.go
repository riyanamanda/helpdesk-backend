package sliceutil

func Map[T, R any](items []T, fn func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}
