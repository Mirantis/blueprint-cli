package utils

// filter applies the given function f to each element of the slice and returns
// a new slice containing only the elements for which f returns true. TODO:

func Filter[T comparable](slice []T, f func(T) bool) []T {
	var result []T

	for _, value := range slice {
		if f(value) {
			result = append(result, value)
		}
	}

	return result
}
