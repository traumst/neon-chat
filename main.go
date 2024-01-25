package main

import (
	"log"
	"net/http"
	"time"

	"go.chat/controllers"
	"go.chat/utils"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middleware ...Middleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := controllers.SetReqId(r)
		startTime := time.Now()
		log.Printf("--%s-> _req_ %s %s", reqId, r.Method, r.RequestURI)
		rec := utils.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		log.Printf("<-%s-- _res_ %s %s %v, status_code:[%d]",
			reqId, r.Method, r.RequestURI, time.Since(startTime), rec.Status())
	})
}

func ControllerSetup() {
	middleware := []Middleware{LoggerMiddleware}

	http.Handle("/favicon.ico", http.HandlerFunc(controllers.FavIcon))
	http.Handle("/login", http.HandlerFunc(controllers.Login))

	http.Handle("/script/", ChainMiddleware(http.HandlerFunc(controllers.ServeFile), middleware...))

	http.Handle("/message", ChainMiddleware(http.HandlerFunc(controllers.AddMessage), middleware...))
	http.Handle("/message/delete", ChainMiddleware(http.HandlerFunc(controllers.DeleteMessage), middleware...))

	chatController := controllers.ChatController{}
	http.Handle("/chat/poll", http.HandlerFunc(chatController.PollChats))
	http.Handle("/chat/invite", ChainMiddleware(
		http.HandlerFunc(chatController.InviteUser), middleware...))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(chatController.OpenChat), middleware...))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(chatController.AddChat), middleware...))

	http.Handle("/", ChainMiddleware(http.HandlerFunc(controllers.Home), middleware...))
}

func main() {
	log.Println("Setting up log middleware")
	log.Println("Setting up controllers")
	ControllerSetup()
	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
