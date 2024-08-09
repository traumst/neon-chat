package controller

import (
	"log"
	"net/http"

	"prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/state"
	h "prplchat/src/utils/http"
)

func Setup(state *state.State, db *db.DBConn) {
	// loaded in reverse order
	allMiddleware := []Middleware{LoggerMiddleware, ReqIdMiddleware}

	handleAvatar(state, db, allMiddleware)
	handleAuth(state, db, allMiddleware)
	handleUser(state, db, allMiddleware)
	handleChat(state, db, allMiddleware)
	handleMsgs(state, db, allMiddleware)
	handleSettings(state, db, allMiddleware)
	handleStaticFiles()

	// live updates
	http.Handle("/poll", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PollUpdates(state, db, w, r)
		}), allMiddleware))

	// home, default
	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := handler.ReadSession(state, db, w, r)
			if err != nil || user == nil {
				log.Printf("[%s] home INFO session, %s\n", h.GetReqId(r), err)
				RenderLogin(w, r)
				return
			}
			RenderHome(state, db, w, r, user)
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

func handleSettings(state *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/settings", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			OpenSettings(state, db, w, r)
		}), allMiddleware))
	http.Handle("/settings/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CloseSettings(state, db, w, r)
		}), allMiddleware))
}

func handleMsgs(state *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteMessage(state, db, w, r)
		}), allMiddleware))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddMessage(state, db, w, r)
		}), allMiddleware))
}

func handleChat(state *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/chat/welcome", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Welcome(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CloseChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			OpenChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddChat(state, db, w, r)
		}), allMiddleware))
}

func handleUser(state *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/user/invite", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			InviteUser(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/expel", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ExpelUser(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/leave", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LeaveChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/change", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ChangeUser(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/search", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SearchUsers(state, db, w, r)
		}), allMiddleware))
}

func handleAuth(state *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Login(state, db, w, r)
		}), allMiddleware))
	http.Handle("/signup", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SignUp(state, db, w, r)
		}), allMiddleware))
	http.Handle("/signup-confirm", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ConfirmEmail(state, db, w, r)
		}), allMiddleware))
	http.Handle("/logout", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logout(state, db, w, r)
		}), allMiddleware))
}

func handleAvatar(state *state.State, db *db.DBConn, allMiddleware []Middleware) {
	http.Handle("/avatar/add", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddAvatar(state, db, w, r)
		}), allMiddleware))
	http.Handle("/avatar", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GetAvatar(state, db, w, r)
		}), allMiddleware))
}
