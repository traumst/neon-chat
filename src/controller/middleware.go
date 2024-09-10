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
			user, _ := handler.ReadSession(state, db, w, r)
			ctx := r.Context()
			if user != nil {
				ctx = context.WithValue(ctx, utils.ActiveUser, user)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx.Value(utils.ActiveUser) == nil {
			w.WriteHeader(http.StatusUnauthorized)
			http.Header.Add(w.Header(), "HX-Refresh", "true")
			w.Write([]byte("unauthorized"))
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func StampMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := h.SetReqId(r, nil)
		log.Printf("[%s] StampMiddleware BEGIN", reqId)

		ctx := context.WithValue(r.Context(), utils.ReqIdKey, reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Printf("[%s] StampMiddleware END", reqId)
	})
}

func StatefulWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Context().Value(utils.ReqIdKey).(string)
		log.Printf("[%s] LoggerMiddleware BEGIN %s %s", reqId, r.Method, r.RequestURI)
		startTime := time.Now()
		rec := h.StatefulWriter{ResponseWriter: w}

		next.ServeHTTP(&rec, r)
		log.Printf("[%s] LoggerMiddleware END %s %s status_code:[%d] in %v",
			reqId,
			r.Method,
			r.RequestURI,
			rec.Status(),
			time.Since(startTime))
	})
}

func AppStateMiddleware(state *state.State) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), utils.AppState, state)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func DBConnMiddleware(db *db.DBConn) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), utils.DBConn, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// func TransactionMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		db := r.Context().Value(utils.DBConn).(*db.DBConn)
// 		dbTx, err := db.AddAuth()
// 		ctx := context.WithValue(r.Context(), utils.DBConn, dbTx)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
