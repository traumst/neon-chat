package middleware

import (
	"context"
	"log"
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"net/http"
)

func DBConnMiddleware(db *db.DBConn) Middleware {
	return Middleware{
		Name: "DBConn",
		Func: func(next http.Handler) http.Handler {
			log.Println("TRACE with db conn middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println(r.Context().Value(consts.ReqIdKey).(string), "TRACE attaching db conn to request context")
				ctx := context.WithValue(r.Context(), consts.DBConn, db)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	}
}
