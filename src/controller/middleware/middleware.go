package middleware

import "net/http"

type Middleware struct {
	Name string
	Func func(http.Handler) http.Handler
}
