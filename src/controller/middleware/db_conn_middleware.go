package middleware

import (
	"context"
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/db"
)

func DBConnMiddleware(dbConn *db.DBConn) Middleware {
	return Middleware{
		Name: "DBConn",
		Func: func(next http.Handler) http.Handler {
			//log.Println("TRACE with db conn middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// log.Printf("TRACE [%s] attaching dbConn to ctx\n", r.Context().Value(consts.ReqIdKey).(string))
				ctx := context.WithValue(r.Context(), consts.DBConn, dbConn)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	}
}
