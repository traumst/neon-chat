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
	"go.chat/model/net"
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

func StatefulWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := net.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&writer, r)
	})
}

func DBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--%s-> BEGIN %s %s", utils.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := net.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		log.Printf("<-%s-- END %s %s status_code:[%d] in %v",
			utils.GetReqId(r),
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}

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
	http.Handle("/chat/user", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DropUser(app, w, r)
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
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	// write log to both file and stderr
	multi := io.MultiWriter(logFile, os.Stderr)
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
	dbPath := "db/chat.db"
	db, err := db.ConnectDB(dbPath)
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
