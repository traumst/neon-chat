package middleware

import (
	"log"
	"neon-chat/src/consts"
	h "neon-chat/src/utils/http"
	"net/http"
	"time"
)

func StatefulWriterMiddleware() Middleware {
	return Middleware{
		Name: "StatefulWriter",
		Func: func(next http.Handler) http.Handler {
			//log.Println("TRACE with auth read middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqId := r.Context().Value(consts.ReqIdKey).(string)
				//log.Printf("TRACE IN [%s] add stateful writer to '%s' '%s'\n", reqId, r.Method, r.RequestURI)
				startTime := time.Now()
				rec := h.StatefulWriter{ResponseWriter: w}

				next.ServeHTTP(&rec, r)
				log.Printf("TRACE OUT [%s] '%s' '%s' with status_code:[%d] in %v\n",
					reqId,
					r.Method,
					r.RequestURI,
					rec.Status(),
					time.Since(startTime))
			})
		}}
}
