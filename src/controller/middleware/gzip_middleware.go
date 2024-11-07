package middleware

import (
	"compress/gzip"
	"io"
	"log"
	h "neon-chat/src/utils/http"
	"net/http"
	"strings"
)

func GZipMiddleware() Middleware {
	return Middleware{
		Name: "Gzip",
		Func: func(next http.Handler) http.Handler {
			log.Println("TRACE with gzip middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
					log.Println("TRACE gzip NO", r.Header.Get("Accept-Encoding"))
					next.ServeHTTP(w, r)
					return
				}
				log.Println("TRACE gzip YES", r.Header.Get("Accept-Encoding"))
				h.SetGzipHeaders(&w)
				gz := gzip.NewWriter(w)
				defer gz.Close()

				gzipWriter := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
				next.ServeHTTP(gzipWriter, r)
			})
		}}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
