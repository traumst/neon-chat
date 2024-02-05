package main

import (
	"log"
	"net/http"
	"time"

	"go.chat/controller"
	"go.chat/utils"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middleware []Middleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SetReqId(r)
		log.Printf("--%s-> _req_ %s %s", utils.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := utils.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		log.Printf("<-%s-- _res_ %s %s status_code:[%d] in %v",
			utils.GetReqId(r), r.Method, r.RequestURI, rec.Status(), time.Since(startTime))
	})
}

func ControllerSetup() {
	noLog := []Middleware{}
	allMiddleware := []Middleware{LoggerMiddleware}

	http.Handle("/favicon.ico", ChainMiddleware(http.HandlerFunc(controller.FavIcon), noLog))
	http.Handle("/login", ChainMiddleware(http.HandlerFunc(controller.Login), allMiddleware))
	http.Handle("/script/", ChainMiddleware(http.HandlerFunc(controller.ServeFile), allMiddleware))

	http.Handle("/message", ChainMiddleware(http.HandlerFunc(controller.AddMessage), allMiddleware))
	http.Handle("/message/delete", ChainMiddleware(http.HandlerFunc(controller.DeleteMessage), allMiddleware))

	chatController := controller.ChatController{}
	http.Handle("/chat/poll", ChainMiddleware(http.HandlerFunc(chatController.PollUpdates), noLog))
	http.Handle("/chat/invite", ChainMiddleware(
		http.HandlerFunc(chatController.InviteUser), allMiddleware))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(chatController.OpenChat), allMiddleware))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(chatController.AddChat), allMiddleware))

	http.Handle("/", ChainMiddleware(http.HandlerFunc(controller.Home), allMiddleware))
}

func main() {
	log.Println("Setting up log middleware")
	log.Println("Setting up controllers")
	ControllerSetup()
	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
