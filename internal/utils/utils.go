package utils

func Coalesce[T any](val *T, defaultVal T) T {
	if val != nil {
		return *val
	}
	return defaultVal
} 