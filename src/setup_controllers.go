package src

import (
	"net/http"

	"neon-chat/src/controller"
	"neon-chat/src/controller/middleware"
	"neon-chat/src/db"
	"neon-chat/src/state"
	"neon-chat/src/utils/config"
)

func SetupControllers(state *state.State, dbConn *db.DBConn, limit config.RpsLimit) {
	withTx := middleware.TransactionMiddleware()
	authValidate := middleware.AuthValidateMiddleware()
	conn := middleware.DBConnMiddleware(dbConn)
	appState := middleware.AppStateMiddleware(state)
	authRead := middleware.AuthReadMiddleware(state, dbConn)
	writer := middleware.StatefulWriterMiddleware()
	gzip := middleware.GZipMiddleware()
	stamp := middleware.StampMiddleware()
	accessCtrl := middleware.AccessControlMiddleware()
	throttleTotal := middleware.ThrottlingTotalMiddleware(limit.TotalRPS, limit.TotalBurst)
	throttleUser := middleware.ThrottlingUserMiddleware(limit.UserRPS, limit.UserBurst)
	maintenance := middleware.MaintenanceMiddleware()
	// TODO uncomment for live
	//recovery := middleware.RecoveryMiddleware()

	// middlewares are loaded in reverse order - lifo
	minMiddleware := middleware.Middlewares{
		conn,
		appState,
		authRead,
		writer,
		gzip,
		stamp,
		accessCtrl,
		throttleTotal,
		throttleUser,
		maintenance,
		//recovery,
	}
	// sse hates gzip
	maxUnzippedMiddleware := middleware.Middlewares{
		withTx,
		authValidate,
		conn,
		appState,
		authRead,
		writer,
		stamp,
		accessCtrl,
		throttleTotal,
		throttleUser,
		maintenance,
		//recovery,
	}
	// most endpoints need all middlewares
	maxMiddleware := middleware.Middlewares{
		withTx,
		authValidate,
		conn,
		appState,
		authRead,
		writer,
		gzip,
		stamp,
		accessCtrl,
		throttleTotal,
		throttleUser,
		maintenance,
		//recovery,
	}

	handleStaticFiles(minMiddleware)
	handleAuth(minMiddleware)

	handleContacts(maxMiddleware)
	handleAvatar(maxMiddleware)
	handleUser(maxMiddleware)
	handleChat(maxMiddleware)
	handleMsgs(maxMiddleware)
	handleSettings(maxMiddleware)

	// live updates
	//log.Println("...live update polling middleware", maxMiddleware)
	http.Handle("/poll", maxUnzippedMiddleware.Chain(
		http.HandlerFunc(controller.PollUpdates)))

	// home, default
	//log.Println("...home middleware", minMiddleware)
	http.Handle("/", minMiddleware.Chain(
		http.HandlerFunc(controller.NavigateHome)))
}

func handleContacts(mw middleware.Middlewares) {
	//log.Println("...contacts middleware", mw)
	http.Handle("/infocard", mw.Chain(
		http.HandlerFunc(controller.GetUserInfo)))
}

func handleStaticFiles(mw middleware.Middlewares) {
	//log.Println("...static files middleware", mw)
	http.Handle("/favicon.ico", mw.Chain(
		http.HandlerFunc(controller.FavIcon)))
	http.Handle("/favicon.svg", mw.Chain(
		http.HandlerFunc(controller.FavIcon)))
	http.Handle("/icon/", mw.Chain(
		http.HandlerFunc(controller.ServeFile)))
	http.Handle("/script/", mw.Chain(
		http.HandlerFunc(controller.ServeFile)))
	http.Handle("/css/", mw.Chain(
		http.HandlerFunc(controller.ServeFile)))
}

func handleSettings(mw middleware.Middlewares) {
	//log.Println("...settings middleware", mw)
	http.Handle("/settings", mw.Chain(
		http.HandlerFunc(controller.OpenSettings)))
	http.Handle("/settings/close", mw.Chain(
		http.HandlerFunc(controller.CloseSettings)))
}

func handleMsgs(mw middleware.Middlewares) {
	//log.Println("...message middleware", mw)
	http.Handle("/message/delete", mw.Chain(
		http.HandlerFunc(controller.DeleteMessage)))
	http.Handle("/message", mw.Chain(
		http.HandlerFunc(controller.AddMessage)))
	http.Handle("/message/quote", mw.Chain(
		http.HandlerFunc(controller.QuoteMessage)))
}

func handleChat(mw middleware.Middlewares) {
	//log.Println("...chat middleware", mw)
	http.Handle("/chat/welcome", mw.Chain(
		http.HandlerFunc(controller.Welcome)))
	http.Handle("/chat/delete", mw.Chain(
		http.HandlerFunc(controller.DeleteChat)))
	http.Handle("/chat/close", mw.Chain(
		http.HandlerFunc(controller.CloseChat)))
	http.Handle("/chat/", mw.Chain(
		http.HandlerFunc(controller.OpenChat)))
	http.Handle("/chat", mw.Chain(
		http.HandlerFunc(controller.AddChat)))
}

func handleUser(mw middleware.Middlewares) {
	//log.Println("...user middleware", mw)
	http.Handle("/user/invite", mw.Chain(
		http.HandlerFunc(controller.InviteUser)))
	http.Handle("/user/expel", mw.Chain(
		http.HandlerFunc(controller.ExpelUser)))
	http.Handle("/user/leave", mw.Chain(
		http.HandlerFunc(controller.LeaveChat)))
	http.Handle("/user/change", mw.Chain(
		http.HandlerFunc(controller.ChangeUser)))
	http.Handle("/user/search", mw.Chain(
		http.HandlerFunc(controller.SearchUsers)))
}

func handleAuth(mw middleware.Middlewares) {
	//log.Println("...auth middleware", mw)
	http.Handle("/login", mw.Chain(
		http.HandlerFunc(controller.Login)))
	http.Handle("/signup", mw.Chain(
		http.HandlerFunc(controller.SignUp)))
	http.Handle("/signup-confirm", mw.Chain(
		http.HandlerFunc(controller.ConfirmEmail)))
	http.Handle("/logout", mw.Chain(
		http.HandlerFunc(controller.Logout)))
}

func handleAvatar(mw middleware.Middlewares) {
	//log.Println("...avatar middleware", mw)
	http.Handle("/avatar/add", mw.Chain(
		http.HandlerFunc(controller.AddAvatar)))
	http.Handle("/avatar", mw.Chain(
		http.HandlerFunc(controller.GetAvatar)))
}
