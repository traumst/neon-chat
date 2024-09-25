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
					log.Println(reqId, "TRACE checking for panic")
					if err := recover(); err != nil {
						log.Printf(reqId, "FATAL recovered from panic: %v", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
					log.Println(reqId, "TRACE request recovered")
				}(r.Context().Value(consts.ReqIdKey).(string))
				next.ServeHTTP(w, r)
			})
		}}
}
