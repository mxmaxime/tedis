package utils

import "strconv"

func StringOr(val1 string, val2 string) string {
	if val1 == "" {
		return val2
	}
	return val1
}

func StringOrInt(val1 string, val2 int) int {
	if val1 == "" {
		intVal, err := strconv.Atoi(val1)
		if err != nil {
			return val2
		}
		return intVal
	}

	return val2
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
