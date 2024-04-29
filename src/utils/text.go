package utils

import "strings"

// replaces \s\s, \t and \n to single \s
func Trim(s string) string {
	res := strings.ReplaceAll(s, "\n", " ")
	// remove double spaces
	for strings.Contains(res, "  ") {
		res = strings.ReplaceAll(res, "  ", " ")
	}
	return res
}

func SanitizeHTML(s string) {
	panic("not implemented")
}

func SanitizeJS(s string) {
	panic("not implemented")
}
