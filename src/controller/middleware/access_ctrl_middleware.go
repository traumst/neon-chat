package middleware

import (
	h "neon-chat/src/utils/http"
	"net/http"
)

func AccessControlMiddleware() Middleware {
	return Middleware{
		Name: "AccessCtrl",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.SetAccessControlHeaders(&w)
				next.ServeHTTP(w, r)
			})
		}}
}
