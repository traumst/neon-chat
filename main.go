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
	"go.chat/db"
	"go.chat/model"
	"go.chat/utils"
)

func ControllerSetup(app *model.AppState, conn *db.DBConn) {
	// static files
	http.Handle("/favicon.ico", http.HandlerFunc(controller.FavIcon))
	http.Handle("/script/", http.HandlerFunc(controller.ServeFile))
	// TEMP
	http.Handle("/gen2", ChainMiddleware(
		http.HandlerFunc(controller.Gen2),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// sessions
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Login(app, conn, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/logout", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Logout(app, conn, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// chat
	http.Handle("/chat/invite", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.InviteUser(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// live updates
	http.Handle("/poll", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.PollUpdates(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// message
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteMessage(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddMessage(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// home, default
	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Home(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
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
	// write log to both file and stderr
	multi := io.MultiWriter(file, os.Stderr)
	log.SetOutput(multi)
	// parse args
	args, err := utils.ArgsRead()
	if err != nil {
		log.Printf("Error parsing args: %v\n", err)
		log.Println(utils.ArgsHelp())
		os.Exit(13)
	}
	log.Printf("  args: %v\n", *args)
	// TODO args.DBPath
	db, err := db.ConnectDB("db/chat.db")
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	log.Println("Setting up application")
	app := &model.ApplicationState
	log.Println("Setting up controllers")
	ControllerSetup(app, db)
	log.Printf("Starting server at port [%d]\n", args.Port)
	runtineErr := http.ListenAndServe(fmt.Sprintf(":%d", args.Port), nil)
	log.Fatal(runtineErr)
}
