package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
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
	//log.Printf("---> ReqIdMiddleware TRACE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("---> ReqIdMiddleware TRACE reqId set to [%s]", utils.SetReqId(r, nil))
		next.ServeHTTP(w, r)
	})
}

func DBMiddleware(next http.Handler) http.Handler {
	log.Printf("---> DBMiddleware TRACE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--%s-> LoggerMiddleware BEGIN %s %s", utils.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := utils.StatefulWriter{ResponseWriter: w}

		next.ServeHTTP(&rec, r)
		log.Printf("<-%s-- LoggerMiddleware END %s %s status_code:[%d] in %v",
			utils.GetReqId(r),
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}

func MinificationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		m.AddFunc("text/html", html.Minify)
		m.AddFunc("text/javascript", js.Minify)

		log.Printf("--%s-> MinificationMiddleware TRACE", utils.GetReqId(r))
		// 1
		//m.Middleware(next).ServeHTTP(w, r)
		// 2
		// mw := m.ResponseWriter(w, r)
		// defer mw.Close()
		// next.ServeHTTP(mw, r)
		// 3
		mw := m.ResponseWriter(w, r)
		defer mw.Close()
		next.ServeHTTP(mw, r)
	})
}
