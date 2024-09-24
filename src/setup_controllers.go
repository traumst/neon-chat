package src

import (
	"net/http"

	"neon-chat/src/controller"
	"neon-chat/src/controller/middleware"
	"neon-chat/src/db"
	"neon-chat/src/state"
)

func SetupControllers(state *state.State, dbConn *db.DBConn) {
	withTx := middleware.TransactionMiddleware()
	authValidate := middleware.AuthValidateMiddleware()
	conn := middleware.DBConnMiddleware(dbConn)
	appState := middleware.AppStateMiddleware(state)
	authRead := middleware.AuthReadMiddleware(state, dbConn)
	writer := middleware.StatefulWriterMiddleware()
	stamp := middleware.StampMiddleware()
	//recovery := middleware.RecoveryMiddleware()

	// middlewares are loaded in reverse order - lifo
	minMiddlewareSet := []middleware.Middleware{
		conn,
		appState,
		authRead,
		writer,
		stamp,
		//recovery,
	}
	// middlewares are loaded in reverse order - lifo
	maxMiddleware := []middleware.Middleware{
		withTx,
		authValidate,
		conn,
		appState,
		authRead,
		writer,
		stamp,
		//recovery,
	}

	handleStaticFiles(minMiddlewareSet)
	handleAuth(minMiddlewareSet)

	handleContacts(maxMiddleware)
	handleAvatar(maxMiddleware)
	handleUser(maxMiddleware)
	handleChat(maxMiddleware)
	handleMsgs(maxMiddleware)
	handleSettings(maxMiddleware)

	// live updates
	http.Handle("/poll", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.PollUpdates(state, dbConn, w, r)
		}), maxMiddleware))

	// home, default
	http.Handle("/", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.NavigateHome(w, r)
		}), minMiddlewareSet))
}

func handleContacts(mw []middleware.Middleware) {
	http.Handle("/infocard", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.GetUserInfo(w, r)
		}), mw))
}

func handleStaticFiles(mw []middleware.Middleware) {
	http.Handle("/favicon.ico", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		}), mw))
	http.Handle("/favicon.svg", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		}), mw))
	http.Handle("/icon/", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), mw))
	http.Handle("/script/", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), mw))
	http.Handle("/css/", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		}), mw))
}

func handleSettings(mw []middleware.Middleware) {
	http.Handle("/settings", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenSettings(w, r)
		}), mw))
	http.Handle("/settings/close", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseSettings(w, r)
		}), mw))
}

func handleMsgs(mw []middleware.Middleware) {
	http.Handle("/message/delete", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteMessage(w, r)
		}), mw))
	http.Handle("/message", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddMessage(w, r)
		}), mw))
	http.Handle("/message/quote", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.QuoteMessage(w, r)
		}), mw))
}

func handleChat(mw []middleware.Middleware) {
	http.Handle("/chat/welcome", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Welcome(w, r)
		}), mw))
	http.Handle("/chat/delete", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteChat(w, r)
		}), mw))
	http.Handle("/chat/close", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseChat(w, r)
		}), mw))
	http.Handle("/chat/", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenChat(w, r)
		}), mw))
	http.Handle("/chat", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddChat(w, r)
		}), mw))
}

func handleUser(mw []middleware.Middleware) {
	http.Handle("/user/invite", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.InviteUser(w, r)
		}), mw))
	http.Handle("/user/expel", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ExpelUser(w, r)
		}), mw))
	http.Handle("/user/leave", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.LeaveChat(w, r)
		}), mw))
	http.Handle("/user/change", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ChangeUser(w, r)
		}), mw))
	http.Handle("/user/search", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SearchUsers(w, r)
		}), mw))
}

func handleAuth(mw []middleware.Middleware) {
	http.Handle("/login", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Login(w, r)
		}), mw))
	http.Handle("/signup", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SignUp(w, r)
		}), mw))
	http.Handle("/signup-confirm", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ConfirmEmail(w, r)
		}), mw))
	http.Handle("/logout", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Logout(w, r)
		}), mw))
}

func handleAvatar(mw []middleware.Middleware) {
	http.Handle("/avatar/add", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddAvatar(w, r)
		}), mw))
	http.Handle("/avatar", middleware.ChainMiddlewares(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.GetAvatar(w, r)
		}), mw))
}
