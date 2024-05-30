package controller

import (
	"log"
	"net/http"
	"time"

	h "prplchat/src/utils/http"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middleware []Middleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func ReqIdMiddleware(next http.Handler) http.Handler {
	//log.Printf("ReqIdMiddleware TRACE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = h.SetReqId(r, nil)
		//log.Printf("ReqIdMiddleware TRACE reqId set to [%s]", reqId)
		next.ServeHTTP(w, r)
	})
}

func DBMiddleware(next http.Handler) http.Handler {
	log.Printf("DBMiddleware TRACE")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		next.ServeHTTP(w, r)
	})
}

// func AuthMiddleware(next http.HandlerFunc) http.Handler {
// 	log.Printf("AuthMiddleware TRACE")
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		rec := utils.StatefulWriter{ResponseWriter: w}
// 		user, err := handler.ReadSession(app, w, r)
// 		if err != nil || user == nil {
// 			log.Printf("[%s] AuthMiddleware WARN user, %s\n", utils.GetReqId(r), err)
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 			w.Write([]byte("User is unauthorized"))
// 			return
// 		}
// 		next(w, r)
// 	})
// }

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] LoggerMiddleware BEGIN %s %s", h.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := h.StatefulWriter{ResponseWriter: w}

		next.ServeHTTP(&rec, r)
		log.Printf("[%s] LoggerMiddleware END %s %s status_code:[%d] in %v",
			h.GetReqId(r),
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}
