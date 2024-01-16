package utils

import "net/http"

type StatefulWriter struct {
	http.ResponseWriter
	status int
}

func (rec *StatefulWriter) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *StatefulWriter) Status() int {
	return rec.status
}
