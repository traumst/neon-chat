package middleware

import (
	"log"
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/db"
	h "neon-chat/src/utils/http"
)

// TODO test
func DBMaintenanceMiddleware(dbConn *db.DBConn) Middleware {
	return Middleware{
		Name: "DBMaintenance",
		Func: func(next http.Handler) http.Handler {
			//log.Println("TRACE with db maintenance middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// log.Printf("TRACE [%s] checking dbConn is not under maintenance\n", r.Context().Value(consts.ReqIdKey).(string))
				if !dbConn.IsAvailable(maxWait) {
					log.Printf("WARN [%s] DB is still unavailable after %s\n", r.Context().Value(consts.ReqIdKey).(string), maxWait)
					h.SetRetryAfterHeader(&w, retryAfter)
					http.Error(w, "Back in a jiff", http.StatusServiceUnavailable)
					return
				}

				next.ServeHTTP(w, r)
			})
		},
	}
}
