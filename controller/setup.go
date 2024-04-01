package controller

import (
	"net/http"

	"go.chat/db"
	"go.chat/model"
)

func Setup(app *model.AppState, conn *db.DBConn, loadLocal bool) {
	// loaded in reverse order
	allMiddleware := []Middleware{
		MinificationMiddleware,
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
	http.Handle("/user/invite", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LeaveChat(app, w, r)
		}), allMiddleware))
}

func handleAuth(app *model.AppState, conn *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Login(app, conn, w, r)
		}), allMiddleware))
	http.Handle("/logout", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logout(w, r)
		}), allMiddleware))
}
