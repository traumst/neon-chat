package middleware

import (
	"log"
	"neon-chat/src/consts"
	"net/http"
)

func RecoveryMiddleware() Middleware {
	return Middleware{
		Name: "Recovery",
		Func: func(next http.Handler) http.Handler {
			//log.Println("TRACE with recovery middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func(reqId string) {
					log.Printf("TRACE [%s] checking for panic '%s' '%s'\n", reqId, r.Method, r.RequestURI)
					if err := recover(); err != nil {
						log.Printf("FATAL [%s] recovered from panic: %v", reqId, err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				}(r.Context().Value(consts.ReqIdKey).(string))
				next.ServeHTTP(w, r)
			})
		}}
}
