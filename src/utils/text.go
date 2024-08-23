package utils

import "strings"

func Contains[T comparable](arr []T, item T) bool {
	for _, a := range arr {
		if a == item {
			return true
		}
	}
	return false
}

// applies all other trims
func SanitizeInput(s string) string {
	s = ReplaceWithSingleSpace(s)
	s = RemoveSpecialChars(s)
	s = strings.Trim(s, " ")
	return s
}

// replaces \s\s, \t and \n to single \s
func ReplaceWithSingleSpace(s string) string {
	if len(s) <= 0 {
		return ""
	}
	res := strings.ReplaceAll(s, "\n", " ")
	res = strings.ReplaceAll(res, "\t", " ")
	// remove double spaces
	for strings.Contains(res, "  ") {
		res = strings.ReplaceAll(res, "  ", " ")
	}
	return res
}

func RemoveSpecialChars(s string) string {
	if len(s) <= 0 {
		return ""
	}
	res := strings.ReplaceAll(s, "\\", "")
	res = strings.ReplaceAll(res, "\"", "")
	res = strings.ReplaceAll(res, "'", "")
	res = strings.ReplaceAll(res, "<", "")
	res = strings.ReplaceAll(res, ">", "")
	res = strings.ReplaceAll(res, "/", "")
	res = strings.ReplaceAll(res, "[", "")
	res = strings.ReplaceAll(res, "]", "")
	return res
}
