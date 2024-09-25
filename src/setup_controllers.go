package src

import (
	"log"
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
	// TODO uncomment for live
	//recovery := middleware.RecoveryMiddleware()

	// middlewares are loaded in reverse order - lifo
	minMiddlewareSet := middleware.Middlewares{
		conn,
		appState,
		authRead,
		writer,
		stamp,
		//recovery,
	}
	// middlewares are loaded in reverse order - lifo
	maxMiddleware := middleware.Middlewares{
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
	log.Println("...live update polling middleware", maxMiddleware)
	http.Handle("/poll", maxMiddleware.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.PollUpdates(state, dbConn, w, r)
		})))

	// home, default
	log.Println("...home middleware", minMiddlewareSet)
	http.Handle("/", minMiddlewareSet.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.NavigateHome(w, r)
		})))
}

func handleContacts(mw middleware.Middlewares) {
	log.Println("...contacts middleware", mw)
	http.Handle("/infocard", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.GetUserInfo(w, r)
		})))
}

func handleStaticFiles(mw middleware.Middlewares) {
	log.Println("...static files middleware", mw)
	http.Handle("/favicon.ico", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		})))
	http.Handle("/favicon.svg", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.FavIcon(w, r)
		})))
	http.Handle("/icon/", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		})))
	http.Handle("/script/", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		})))
	http.Handle("/css/", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ServeFile(w, r)
		})))
}

func handleSettings(mw middleware.Middlewares) {
	log.Println("...settings middleware", mw)
	http.Handle("/settings", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenSettings(w, r)
		})))
	http.Handle("/settings/close", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseSettings(w, r)
		})))
}

func handleMsgs(mw middleware.Middlewares) {
	log.Println("...message middleware", mw)
	http.Handle("/message/delete", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteMessage(w, r)
		})))
	http.Handle("/message", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddMessage(w, r)
		})))
	http.Handle("/message/quote", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.QuoteMessage(w, r)
		})))
}

func handleChat(mw middleware.Middlewares) {
	log.Println("...chat middleware", mw)
	http.Handle("/chat/welcome", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Welcome(w, r)
		})))
	http.Handle("/chat/delete", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.DeleteChat(w, r)
		})))
	http.Handle("/chat/close", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.CloseChat(w, r)
		})))
	http.Handle("/chat/", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.OpenChat(w, r)
		})))
	http.Handle("/chat", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddChat(w, r)
		})))
}

func handleUser(mw middleware.Middlewares) {
	log.Println("...user middleware", mw)
	http.Handle("/user/invite", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.InviteUser(w, r)
		})))
	http.Handle("/user/expel", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ExpelUser(w, r)
		})))
	http.Handle("/user/leave", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.LeaveChat(w, r)
		})))
	http.Handle("/user/change", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ChangeUser(w, r)
		})))
	http.Handle("/user/search", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SearchUsers(w, r)
		})))
}

func handleAuth(mw middleware.Middlewares) {
	log.Println("...auth middleware", mw)
	http.Handle("/login", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Login(w, r)
		})))
	http.Handle("/signup", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.SignUp(w, r)
		})))
	http.Handle("/signup-confirm", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.ConfirmEmail(w, r)
		})))
	http.Handle("/logout", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.Logout(w, r)
		})))
}

func handleAvatar(mw middleware.Middlewares) {
	log.Println("...avatar middleware", mw)
	http.Handle("/avatar/add", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.AddAvatar(w, r)
		})))
	http.Handle("/avatar", mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controller.GetAvatar(w, r)
		})))
}
