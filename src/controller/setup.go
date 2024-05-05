package controller

import (
	"net/http"

	"go.chat/src/db"
	"go.chat/src/handler"
)

func Setup(app *handler.AppState, conn *db.DBConn) {
	// loaded in reverse order
	allMiddleware := []Middleware{LoggerMiddleware, ReqIdMiddleware}

	handleAuth(app, conn, allMiddleware)
	handleUser(app, conn, allMiddleware)
	handleChat(app, allMiddleware)
	handleMsgs(app, allMiddleware)
	handleSettings(app, allMiddleware)
	handleStaticFiles()

	// live updates
	http.Handle("/poll", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PollUpdates(app, w, r)
		}), allMiddleware))

	// home, default
	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderHome(app, w, r)
		}), allMiddleware))
}

func handleStaticFiles() {
	// loaded in reverse order
	minMiddleware := []Middleware{ReqIdMiddleware}

	http.Handle("/favicon.ico", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			FavIcon(w, r)
		}), minMiddleware))
	http.Handle("/icon/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeFile(w, r)
		}), minMiddleware))
	http.Handle("/script/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeFile(w, r)
		}), minMiddleware))
	http.Handle("/css/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeFile(w, r)
		}), minMiddleware))
}

func handleSettings(app *handler.AppState, allMiddleware []Middleware) {
	http.Handle("/settings", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			OpenSettings(app, w, r)
		}), allMiddleware))
	http.Handle("/settings/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CloseSettings(app, w, r)
		}), allMiddleware))
}

func handleMsgs(app *handler.AppState, allMiddleware []Middleware) {
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteMessage(app, w, r)
		}), allMiddleware))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddMessage(app, w, r)
		}), allMiddleware))
}

func handleChat(app *handler.AppState, allMiddleware []Middleware) {
	http.Handle("/chat/welcome", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Welcome(app, w, r)
		}), allMiddleware))
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

func handleUser(app *handler.AppState, conn *db.DBConn, allMiddleware []Middleware) {
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
	http.Handle("/user/change", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ChangeUser(app, w, r)
		}), allMiddleware))
}

func handleAuth(app *handler.AppState, conn *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Login(app, conn, w, r)
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
