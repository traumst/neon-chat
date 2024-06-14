package controller

import (
	"log"
	"net/http"

	"prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/state"
	h "prplchat/src/utils/http"
)

func Setup(app *state.State, db *db.DBConn) {
	// loaded in reverse order
	allMiddleware := []Middleware{LoggerMiddleware, ReqIdMiddleware}

	handleAvatar(app, db, allMiddleware)
	handleAuth(app, db, allMiddleware)
	handleUser(app, db, allMiddleware)
	handleChat(app, db, allMiddleware)
	handleMsgs(app, db, allMiddleware)
	handleSettings(app, db, allMiddleware)
	handleStaticFiles()

	// live updates
	http.Handle("/poll", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PollUpdates(app, db, w, r)
		}), allMiddleware))

	// home, default
	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := handler.ReadSession(app, db, w, r)
			if err != nil || user == nil {
				log.Printf("[%s] home INFO session, %s\n", h.GetReqId(r), err)
				RenderLogin(w, r)
				return
			}
			RenderHome(app, db, w, r, user)
		}), allMiddleware))
}

func handleStaticFiles() {
	// loaded in reverse order
	minMiddleware := []Middleware{ReqIdMiddleware}

	http.Handle("/favicon.ico", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			FavIcon(w, r)
		}), minMiddleware))
	http.Handle("/favicon.svg", ChainMiddleware(
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

func handleSettings(app *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/settings", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			OpenSettings(app, db, w, r)
		}), allMiddleware))
	http.Handle("/settings/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CloseSettings(app, db, w, r)
		}), allMiddleware))
}

func handleMsgs(app *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteMessage(app, db, w, r)
		}), allMiddleware))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddMessage(app, db, w, r)
		}), allMiddleware))
}

func handleChat(app *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/chat/welcome", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Welcome(app, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteChat(app, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CloseChat(app, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			OpenChat(app, db, w, r)
		}), allMiddleware))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddChat(app, db, w, r)
		}), allMiddleware))
}

func handleUser(app *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/user/invite", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			InviteUser(app, db, w, r)
		}), allMiddleware))
	http.Handle("/user/expel", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ExpelUser(app, db, w, r)
		}), allMiddleware))
	http.Handle("/user/leave", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LeaveChat(app, db, w, r)
		}), allMiddleware))
	http.Handle("/user/change", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ChangeUser(app, db, w, r)
		}), allMiddleware))
}

func handleAuth(app *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Login(app, db, w, r)
		}), allMiddleware))
	http.Handle("/signup", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SignUp(app, db, w, r)
		}), allMiddleware))
	http.Handle("/signup-confirm", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ConfirmEmail(app, db, w, r)
		}), allMiddleware))
	http.Handle("/logout", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logout(app, db, w, r)
		}), allMiddleware))
}

func handleAvatar(app *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/avatar/add", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddAvatar(app, db, w, r)
		}), allMiddleware))
	http.Handle("/avatar", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GetAvatar(app, db, w, r)
		}), allMiddleware))
}
