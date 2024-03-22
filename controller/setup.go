package controller

import (
	"log"
	"net/http"
	"time"

	"go.chat/db"
	"go.chat/model"
	"go.chat/model/net"
	"go.chat/utils"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middleware []Middleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func ReqIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SetReqId(r, nil)
		next.ServeHTTP(w, r)
	})
}

func StatefulWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := net.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&writer, r)
	})
}

func DBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--%s-> BEGIN %s %s", utils.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := net.StatefulWriter{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		log.Printf("<-%s-- END %s %s status_code:[%d] in %v",
			utils.GetReqId(r),
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}

func Setup(app *model.AppState, conn *db.DBConn) {
	// static files
	http.Handle("/favicon.ico", http.HandlerFunc(FavIcon))
	http.Handle("/script/", http.HandlerFunc(ServeFile))
	// TEMP
	http.Handle("/gen2", ChainMiddleware(
		http.HandlerFunc(Gen2),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// sessions
	http.Handle("/login", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Login(app, conn, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/logout", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Logout(app, conn, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// chat
	http.Handle("/chat/invite", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			InviteUser(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/expel", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ExpelUser(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/leave", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LeaveChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/close", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			CloseChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			OpenChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/chat", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddChat(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// live updates
	http.Handle("/poll", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PollUpdates(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// message
	http.Handle("/message/delete", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			DeleteMessage(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	http.Handle("/message", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AddMessage(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
	// home, default
	http.Handle("/", ChainMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Home(app, w, r)
		}),
		[]Middleware{LoggerMiddleware, ReqIdMiddleware}))
}
