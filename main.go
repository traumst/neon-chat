package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.chat/controller"
	"go.chat/model"
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
		log.Printf("--%s-> BEGIN %s %s", utils.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := utils.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		log.Printf("<-%s-- END %s %s status_code:[%d] in %v",
			utils.GetReqId(r),
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}

func ControllerSetup(app *model.AppState) {
	noLog := []Middleware{ReqIdMiddleware}
	allMiddleware := []Middleware{LoggerMiddleware}
	// static files
	http.Handle("/favicon.ico", ChainMiddleware(
		http.HandlerFunc(controller.FavIcon),
		noLog))
	http.Handle("/script/", ChainMiddleware(
		http.HandlerFunc(controller.ServeFile),
		noLog))
	// TEMP
	http.Handle("/gen2", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Gen2(w, r)
		}),
		allMiddleware))
	// sessions
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Login(w, r)
		}),
		allMiddleware))
	http.Handle("/logout", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Logout(w, r)
		}),
		allMiddleware))
	// chat
	http.Handle("/chat/invite", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.InviteUser(app, w, r)
		}),
		allMiddleware))
	http.Handle("/chat/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteChat(app, w, r)
		}),
		allMiddleware))
	http.Handle("/chat/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseChat(app, w, r)
		}),
		allMiddleware))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenChat(app, w, r)
		}),
		allMiddleware))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddChat(app, w, r)
		}),
		allMiddleware))
	// live updates
	http.Handle("/poll", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.PollUpdates(app, w, r)
		}),
		allMiddleware))
	// message
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteMessage(app, w, r)
		}),
		allMiddleware))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddMessage(app, w, r)
		}),
		allMiddleware))
	// home, default
	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Home(app, w, r)
		}),
		allMiddleware))
}

func main() {
	// log setup
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	date := strings.Split(timestamp, "T")[0]
	logPath := fmt.Sprintf("log/from-%s.log", date)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// parse args
	args, err := utils.ArgsRead()
	if err != nil {
		log.Printf("Error parsing args: %v\n", err)
		log.Println(utils.ArgsHelp())
		os.Exit(13)
	}
	log.Printf("  args: %v\n", *args)
	log.Println("Setting up application")
	app := &model.ApplicationState
	log.Println("Setting up controllers")
	ControllerSetup(app)
	log.Printf("Starting server at port [%d]\n", args.Port)
	// run server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.Port), nil))
}
