package utils

import "net/http"

type StatefulWriter struct {
	http.ResponseWriter
	status int
	// TODO expand what session can hold
	// userId   uint
	// userName string
}

func StatefulWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&writer, r)
	})
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
