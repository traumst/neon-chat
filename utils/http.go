package utils

import (
	"net/http"
	"strings"
)

func SetReqId(r *http.Request) string {
	reqId := RandStringBytes(5)
	r.Header.Set("X-Request-Id", reqId)
	return reqId
}

func GetReqId(r *http.Request) string {
	return r.Header.Get("X-Request-Id")
}

func SetSseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

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
