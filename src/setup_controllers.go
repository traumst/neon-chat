package src

import (
	"log"
	"net/http"

	"neon-chat/src/controller"
	"neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	h "neon-chat/src/utils/http"
)

func SetupControllers(state *state.State, db *db.DBConn) {
	// loaded in reverse order
	allMiddleware := []controller.Middleware{
		controller.LoggerMiddleware,
		controller.ReqIdMiddleware,
		controller.RecoveryMiddleware,
	}

	handleAvatar(state, db, allMiddleware)
	handleAuth(state, db, allMiddleware)
	handleUser(state, db, allMiddleware)
	handleChat(state, db, allMiddleware)
	handleMsgs(state, db, allMiddleware)
	handleSettings(state, db, allMiddleware)
	handleStaticFiles()

	// live updates
	http.Handle("/poll", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.PollUpdates(state, db, w, r)
		}), allMiddleware))

	// home, default
	http.Handle("/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := handler.ReadSession(state, db, w, r)
			if err != nil || user == nil {
				log.Printf("[%s] home INFO session, %s\n", h.GetReqId(r), err)
				controller.RenderLogin(w, r)
				return
			}
			controller.RenderHome(state, db, w, r, user)
		}), allMiddleware))
}

func handleStaticFiles() {
	// loaded in reverse order
	minMiddleware := []controller.Middleware{controller.ReqIdMiddleware}

	http.Handle("/favicon.ico", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		}), minMiddleware))
	http.Handle("/favicon.svg", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		}), minMiddleware))
	http.Handle("/icon/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), minMiddleware))
	http.Handle("/script/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), minMiddleware))
	http.Handle("/css/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), minMiddleware))
}

func handleSettings(state *state.State, db *db.DBConn, allMiddleware []controller.Middleware) {
	http.Handle("/settings", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenSettings(state, db, w, r)
		}), allMiddleware))
	http.Handle("/settings/close", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseSettings(state, db, w, r)
		}), allMiddleware))
}

func handleMsgs(state *state.State, db *db.DBConn, allMiddleware []controller.Middleware) {
	http.Handle("/message/delete", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteMessage(state, db, w, r)
		}), allMiddleware))
	http.Handle("/message", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddMessage(state, db, w, r)
		}), allMiddleware))
	http.Handle("/message/quote", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.QuoteMessage(state, db, w, r)
		}), allMiddleware))
}

func handleChat(state *state.State, db *db.DBConn, allMiddleware []controller.Middleware) {
	http.Handle("/chat/welcome", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Welcome(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/delete", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/close", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/chat", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddChat(state, db, w, r)
		}), allMiddleware))
}

func handleUser(state *state.State, db *db.DBConn, allMiddleware []controller.Middleware) {
	http.Handle("/user/invite", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.InviteUser(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/expel", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ExpelUser(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/leave", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.LeaveChat(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/change", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ChangeUser(state, db, w, r)
		}), allMiddleware))
	http.Handle("/user/search", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SearchUsers(state, db, w, r)
		}), allMiddleware))
}

func handleAuth(state *state.State, db *db.DBConn, allMiddleware []controller.Middleware) {
	http.Handle("/login", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Login(state, db, w, r)
		}), allMiddleware))
	http.Handle("/signup", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SignUp(state, db, w, r)
		}), allMiddleware))
	http.Handle("/signup-confirm", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ConfirmEmail(state, db, w, r)
		}), allMiddleware))
	http.Handle("/logout", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Logout(state, db, w, r)
		}), allMiddleware))
}

func handleAvatar(state *state.State, db *db.DBConn, allMiddleware []controller.Middleware) {
	http.Handle("/avatar/add", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddAvatar(state, db, w, r)
		}), allMiddleware))
	http.Handle("/avatar", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.GetAvatar(state, db, w, r)
		}), allMiddleware))
}
