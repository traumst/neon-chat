package net

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

func (w *StatefulWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
