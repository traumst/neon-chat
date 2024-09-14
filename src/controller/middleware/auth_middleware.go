package middleware

import (
	"context"
	"log"
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	"net/http"
)

func AuthReadMiddleware(state *state.State, db *db.DBConn) Middleware {
	return Middleware{
		Name: "AuthRead",
		Func: func(next http.Handler) http.Handler {
			log.Println("TRACE with auth read middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println(r.Context().Value(consts.ReqIdKey).(string), "TRACE reading user session auth")
				user, _ := handler.ReadSession(state, db, w, r)
				ctx := r.Context()
				if user != nil {
					ctx = context.WithValue(ctx, consts.ActiveUser, user)
				}

				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}}
}

func AuthValidateMiddleware() Middleware {
	log.Println("TRACE with auth validate middleware")
	return Middleware{
		Name: "AuthValidate",
		Func: func(next http.Handler) http.Handler {
			log.Println("TRACE with auth read middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println(r.Context().Value(consts.ReqIdKey).(string), "TRACE validating user session auth")
				ctx := r.Context()
				if ctx.Value(consts.ActiveUser) == nil {
					w.WriteHeader(http.StatusUnauthorized)
					http.Header.Add(w.Header(), "HX-Refresh", "true")
					w.Write([]byte("unauthorized"))
					return
				}

				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}}
}
