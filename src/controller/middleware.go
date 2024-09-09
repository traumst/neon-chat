package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddlewares(h http.Handler, middleware []Middleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ContextMiddleware(state *state.State, db *db.DBConn) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("ContextMiddleware TRACE IN")

			reqId := h.SetReqId(r, nil)
			ctx := context.WithValue(r.Context(), utils.ReqIdKey, reqId)
			ctx = context.WithValue(ctx, utils.AppState, state)
			ctx = context.WithValue(ctx, utils.DBConn, db)

			user, err := handler.ReadSession(state, db, w, r)
			if err != nil || user == nil {
				log.Printf("ContextMiddleware INFO user has no session, %s\n", err)
			} else {
				ctx = context.WithValue(ctx, utils.ActiveUser, user)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
			log.Printf("ContextMiddleware TRACE OUT")
		})
	}
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] LoggerMiddleware BEGIN %s %s", h.GetReqId(r), r.Method, r.RequestURI)
		startTime := time.Now()
		rec := h.StatefulWriter{ResponseWriter: w}

		next.ServeHTTP(&rec, r)
		log.Printf("[%s] LoggerMiddleware END %s %s status_code:[%d] in %v",
			h.GetReqId(r),
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}
