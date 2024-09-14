package middleware

import (
	"log"
	"neon-chat/src/consts"
	"net/http"
)

type Middleware struct {
	Name string
	Func func(http.Handler) http.Handler
}

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

func GetReqId(r *http.Request) string {
	return r.Context().Value(consts.ReqIdKey).(string)
}
