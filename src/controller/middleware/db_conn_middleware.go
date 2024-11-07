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
				log.Printf("TRACE [%s] attaching dbConn to ctx of '%s' '%s'\n", r.Context().Value(consts.ReqIdKey).(string), r.Method, r.RequestURI)
				ctx := context.WithValue(r.Context(), consts.DBConn, dbConn)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	}
}
