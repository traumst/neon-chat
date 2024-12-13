package utils

func Contains[T comparable](arr []T, item T) bool {
	for _, a := range arr {
		if a == item {
			return true
		}
	}
	return false
}
