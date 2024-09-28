package middleware

import (
	"context"
	"log"
	"neon-chat/src/consts"
	h "neon-chat/src/utils/http"
	"net/http"
)

func StampMiddleware() Middleware {
	return Middleware{
		Name: "Stamp",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqId := h.SetReqId(r, nil)
				log.Printf("TRACE IN [%s] stamp '%s' '%s'\n", reqId, r.Method, r.RequestURI)

				ctx := context.WithValue(r.Context(), consts.ReqIdKey, reqId)
				next.ServeHTTP(w, r.WithContext(ctx))
				log.Printf("TRACE OUT [%s] stamp '%s' '%s'\n", reqId, r.Method, r.RequestURI)
			})
		}}
}
