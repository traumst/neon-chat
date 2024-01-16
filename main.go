package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.chat/controllers"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middleware ...Middleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := RandStringBytes(5)
		r.Header.Set("X-Request-Id", reqId)
		startTime := time.Now()
		log.Printf("--%s-> %s %s", reqId, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("<-%s-- %s %s %v", reqId, r.Method, r.RequestURI, time.Since(startTime))
	})
}

func ControllerSetup() {
	middleware := []Middleware{LoggerMiddleware}

	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(controllers.Home), middleware...))
	http.Handle("/favicon.ico", ChainMiddleware(
		http.HandlerFunc(controllers.FavIcon), middleware...))
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(controllers.Login), middleware...))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(controllers.AddChat), middleware...))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(controllers.OpenChat), middleware...))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(controllers.AddMessage), middleware...))
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(controllers.DeleteMessage), middleware...))
}

func main() {
	log.Println("Setting up log middleware")
	log.Println("Setting up controllers")
	ControllerSetup()
	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
