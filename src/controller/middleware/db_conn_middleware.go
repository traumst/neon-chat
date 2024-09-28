package middleware

import (
	"context"
	"log"
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"net/http"
)

func DBConnMiddleware(dbConn *db.DBConn) Middleware {
	return Middleware{
		Name: "DBConn",
		Func: func(next http.Handler) http.Handler {
			//log.Println("TRACE with db conn middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("TRACE [%s] attaching db conn to request context\n", r.Context().Value(consts.ReqIdKey).(string))
				ctx := context.WithValue(r.Context(), consts.DBConn, dbConn)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	}
}
