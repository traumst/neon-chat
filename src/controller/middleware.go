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

func AuthMiddleware(state *state.State, db *db.DBConn) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := handler.ReadSession(state, db, w, r)
			if err != nil || user == nil {
				log.Printf("AuthMiddleware INFO user in unauthorized, %s\n", err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			ctx := context.WithValue(r.Context(), utils.ActiveUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func StampMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := h.SetReqId(r, nil)
		log.Printf("[%s] StampMiddleware BEGIN", reqId)

		ctx := context.WithValue(r.Context(), utils.ReqIdKey, h.GetReqId(r))
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Printf("[%s] StampMiddleware END", reqId)
	})
}

func StatefulWriterMiddleware(next http.Handler) http.Handler {
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

func DBConnMiddleware(db *db.DBConn) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), utils.DBConn, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AppStateMiddleware(state *state.State) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), utils.AppState, state)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
