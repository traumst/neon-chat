package middleware

import (
	"context"
	"log"
	"neon-chat/src/consts"
	"neon-chat/src/state"
	"net/http"
)

func AppStateMiddleware(state *state.State) Middleware {
	return Middleware{
		Name: "AppState",
		Func: func(next http.Handler) http.Handler {
			log.Println("TRACE with app state middleware")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println(r.Context().Value(consts.ReqIdKey).(string), "TRACE attaching app state to request context")
				ctx := context.WithValue(r.Context(), consts.AppState, state)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}}
}
