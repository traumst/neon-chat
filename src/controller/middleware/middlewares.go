package middleware

import (
	"net/http"
	"strings"
)

type Middlewares []Middleware

func (ms Middlewares) String() string {
	names := make([]string, 0)
	for _, m := range ms {
		names = append(names, m.Name)
	}
	return strings.Join(names, ",")
}

func (mw Middlewares) Chain(h http.Handler) http.Handler {
	count := len(mw)
	names := make([]string, count)
	for i, m := range mw {
		names[i] = m.Name
	}
	//log.Printf("TRACE chaning %d middlewares: %v", count, names)
	for _, m := range mw {
		h = m.Func(h)
	}
	return h
}
