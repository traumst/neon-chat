package utils

import (
	"net/http"
	"strings"
)

func ParseUrlPath(r *http.Request) []string {
	path := strings.Split(r.URL.Path, "/")
	return path
}

func ParseUrlArgs(r *http.Request) map[string]string {
	args := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) == 0 {
			args[k] = v[0]
		}
	}
	return args
}
