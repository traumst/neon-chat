package utils

import "strings"

// replaces \s\s, \t and \n to single \s
func TrimSpaces(s string) string {
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

func TrimSpecial(s string) string {
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
