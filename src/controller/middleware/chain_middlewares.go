package middleware

import (
	"log"
	"net/http"
)

func ChainMiddlewares(h http.Handler, mw []Middleware) http.Handler {
	count := len(mw)
	names := make([]string, count)
	for i, m := range mw {
		names[i] = m.Name
	}
	log.Printf("TRACE chaning %d middlewares: %v", count, names)
	for _, m := range mw {
		h = m.Func(h)
	}
	return h
}
