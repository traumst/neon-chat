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
	ctxMiddleware := controller.ContextMiddleware(state, db)

	//
	minMiddleware := []controller.Middleware{
		controller.ContextMiddleware(state, db),
	}
	handleStaticFiles(minMiddleware)

	// loaded in reverse order - lifo
	allMiddleware := []controller.Middleware{
		ctxMiddleware,
		controller.LoggerMiddleware,
		controller.RecoveryMiddleware,
	}

	handleAvatar(allMiddleware)
	handleAuth(allMiddleware)
	handleUser(allMiddleware)
	handleChat(allMiddleware)
	handleMsgs(allMiddleware)
	handleSettings(allMiddleware)

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

func handleStaticFiles(middleware []controller.Middleware) {
	http.Handle("/favicon.ico", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		}), middleware))
	http.Handle("/favicon.svg", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		}), middleware))
	http.Handle("/icon/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), middleware))
	http.Handle("/script/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), middleware))
	http.Handle("/css/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), middleware))
}

func handleSettings(middleware []controller.Middleware) {
	http.Handle("/settings", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenSettings(w, r)
		}), middleware))
	http.Handle("/settings/close", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseSettings(w, r)
		}), middleware))
}

func handleMsgs(middleware []controller.Middleware) {
	http.Handle("/message/delete", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteMessage(w, r)
		}), middleware))
	http.Handle("/message", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddMessage(w, r)
		}), middleware))
	http.Handle("/message/quote", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.QuoteMessage(w, r)
		}), middleware))
}

func handleChat(middleware []controller.Middleware) {
	http.Handle("/chat/welcome", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Welcome(w, r)
		}), middleware))
	http.Handle("/chat/delete", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteChat(w, r)
		}), middleware))
	http.Handle("/chat/close", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseChat(w, r)
		}), middleware))
	http.Handle("/chat/", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenChat(w, r)
		}), middleware))
	http.Handle("/chat", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddChat(w, r)
		}), middleware))
}

func handleUser(middleware []controller.Middleware) {
	http.Handle("/user/invite", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.InviteUser(w, r)
		}), middleware))
	http.Handle("/user/expel", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ExpelUser(w, r)
		}), middleware))
	http.Handle("/user/leave", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.LeaveChat(w, r)
		}), middleware))
	http.Handle("/user/change", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ChangeUser(w, r)
		}), middleware))
	http.Handle("/user/search", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SearchUsers(w, r)
		}), middleware))
}

func handleAuth(middleware []controller.Middleware) {
	http.Handle("/login", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Login(w, r)
		}), middleware))
	http.Handle("/signup", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SignUp(w, r)
		}), middleware))
	http.Handle("/signup-confirm", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ConfirmEmail(w, r)
		}), middleware))
	http.Handle("/logout", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Logout(w, r)
		}), middleware))
}

func handleAvatar(middleware []controller.Middleware) {
	http.Handle("/avatar/add", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddAvatar(w, r)
		}), middleware))
	http.Handle("/avatar", controller.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.GetAvatar(w, r)
		}), middleware))
}
