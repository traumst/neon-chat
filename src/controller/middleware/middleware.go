package middleware

import (
	"fmt"
	"net/http"
)

type Middleware struct {
	Name string
	Func func(http.Handler) http.Handler
}

func (m Middleware) String() string {
	return fmt.Sprintf("Middleware{Name: %s}", m.Name)
}
