package controller

import (
	"log"
	"net/http"

	"go.chat/db"
	"go.chat/model"
	"go.chat/utils"
)

func Setup(app *model.AppState, conn *db.DBConn, loadLocal bool) {
	// loaded in reverse order
	allMiddleware := []Middleware{
		LoggerMiddleware,
		ReqIdMiddleware}

	handleAuth(app, conn, allMiddleware)
	handleUser(app, conn, allMiddleware)
	handleChat(app, allMiddleware)
	handleMsgs(app, allMiddleware)

	// live updates
	http.Handle("/poll", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PollUpdates(app, w, r)
		}), allMiddleware))

	// home, default
	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Home(app, w, r)
		}), allMiddleware))

	// static files
	http.Handle("/favicon.ico", http.HandlerFunc(FavIcon))
	if loadLocal {
		http.Handle("/script/", http.HandlerFunc(ServeFile))
	}
	http.Handle("/gen2", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Gen2(app, w, r)
		}), allMiddleware))
}

func handleMsgs(app *model.AppState, allMiddleware []Middleware) {
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteMessage(app, w, r)
		}), allMiddleware))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddMessage(app, w, r)
		}), allMiddleware))
}

func handleChat(app *model.AppState, allMiddleware []Middleware) {
	http.Handle("/chat/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteChat(app, w, r)
		}), allMiddleware))
	http.Handle("/chat/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CloseChat(app, w, r)
		}), allMiddleware))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			OpenChat(app, w, r)
		}), allMiddleware))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddChat(app, w, r)
		}), allMiddleware))
}

func handleUser(app *model.AppState, conn *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/user/invite", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			InviteUser(app, conn, w, r)
		}), allMiddleware))
	http.Handle("/user/expel", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ExpelUser(app, w, r)
		}), allMiddleware))
	http.Handle("/user/leave", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LeaveChat(app, w, r)
		}), allMiddleware))
}

func handleAuth(app *model.AppState, conn *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				log.Printf("--%s-> handleAuth.RenderLogin\n", utils.GetReqId(r))
				RenderLogin(app, w, r)
			} else if r.Method == "POST" {
				log.Printf("--%s-> handleAuth.Login\n", utils.GetReqId(r))
				Login(app, conn, w, r)
			} else {
				log.Printf("--%s-> handleAuth invalid method [%s]\n", utils.GetReqId(r), r.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Bad Request"))
				return
			}
		}), allMiddleware))
	http.Handle("/signup", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SignUp(app, conn, w, r)
		}), allMiddleware))
	http.Handle("/logout", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logout(app, w, r)
		}), allMiddleware))
}
