package main

import (
	"log"
	"net/http"
	"time"

	"go.chat/model/net"
	"go.chat/utils"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middleware []Middleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func ReqIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SetReqId(r, nil)
		next.ServeHTTP(w, r)
	})
}

func StatefulWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := net.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&writer, r)
	})
}

func DBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--%s-> BEGIN %s %s", utils.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := net.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		log.Printf("<-%s-- END %s %s status_code:[%d] in %v",
			utils.GetReqId(r),
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}
