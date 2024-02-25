package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func ReqIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SetReqId(r, nil)
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SetReqId(r, nil)
		//log.Printf("--%s-> _req_ %s %s", utils.GetReqId(r), r.Method, r.RequestURI)
		//startTime := time.Now()
		rec := utils.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		//log.Printf("<-%s-- _res_ %s %s status_code:[%d] in %v", utils.GetReqId(r), r.Method, r.RequestURI, rec.Status(), time.Since(startTime))
	})
}

func ControllerSetup() {
	noLog := []Middleware{ReqIdMiddleware}
	allMiddleware := []Middleware{LoggerMiddleware}
	http.Handle("/favicon.ico", ChainMiddleware(http.HandlerFunc(controller.FavIcon), noLog))
	http.Handle("/login", ChainMiddleware(http.HandlerFunc(controller.Login), allMiddleware))
	http.Handle("/script/", ChainMiddleware(http.HandlerFunc(controller.ServeFile), allMiddleware))
	http.Handle("/message", ChainMiddleware(http.HandlerFunc(controller.AddMessage), allMiddleware))
	http.Handle("/message/delete", ChainMiddleware(http.HandlerFunc(controller.DeleteMessage), allMiddleware))
	http.Handle("/poll", ChainMiddleware(http.HandlerFunc(controller.PollUpdates), allMiddleware))

	chatController := controller.ChatController{}
	http.Handle("/chat/invite", ChainMiddleware(http.HandlerFunc(chatController.InviteUser), allMiddleware))
	http.Handle("/chat/", ChainMiddleware(http.HandlerFunc(chatController.OpenChat), allMiddleware))
	http.Handle("/chat", ChainMiddleware(http.HandlerFunc(chatController.AddChat), allMiddleware))

	http.Handle("/gen2", ChainMiddleware(http.HandlerFunc(controller.Gen2), allMiddleware))

	http.Handle("/", ChainMiddleware(http.HandlerFunc(controller.Home), allMiddleware))
}

func main() {
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	date := strings.Split(timestamp, "T")[0]
	logPath := fmt.Sprintf("log/from-%s.log", date)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multi)

	log.Println("Setting up log middleware")
	log.Println("Setting up controllers")
	ControllerSetup()
	log.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
