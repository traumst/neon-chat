package controller

import (
	"context"
	"fmt"
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
				log.Printf("FATAL RecoveryMiddleware recovered from panic: %v", err)
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
		log.Printf("TRACE [%s] StampMiddleware BEGIN", reqId)

		ctx := context.WithValue(r.Context(), utils.ReqIdKey, reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Printf("TRACE [%s] StampMiddleware END", reqId)
	})
}

func StatefulWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Context().Value(utils.ReqIdKey).(string)
		log.Printf("TRACE [%s] StatefulWriterMiddleware BEGIN %s %s", reqId, r.Method, r.RequestURI)
		startTime := time.Now()
		rec := h.StatefulWriter{ResponseWriter: w}

		next.ServeHTTP(&rec, r)
		log.Printf("TRACE [%s] StatefulWriterMiddleware END %s %s status_code:[%d] in %v",
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

func TransactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqId := ctx.Value(utils.ReqIdKey).(string)
		// copy db connection
		db := *ctx.Value(utils.DBConn).(*db.DBConn)
		// populates tx prop
		err := db.OpenTx(reqId)
		if err != nil {
			log.Printf("FATAL [%s] TransactionMiddleware failed to open transaction: %s", reqId, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// ctx with both general conn and per-session tx
		ctx = context.WithValue(ctx, utils.DBConn, &db)

		defer func() {
			if p := recover(); p != nil {
				log.Printf("FATAL [%s] TransactionMiddleware Failed to open transaction: %v", reqId, p)
				db.CloseTx(fmt.Errorf("panic: %v", p), false)
				panic(p) // re-throw the panic after rollback
			} else if code := w.(*h.StatefulWriter).Status(); code >= http.StatusBadRequest {
				db.CloseTx(fmt.Errorf("error status code %d", code), false)
			} else {
				// must explicitly mark when changes are made in BL
				var changesMade bool
				if r.Context().Value(utils.TxChangesKey) != nil {
					changesMade = *r.Context().Value(utils.TxChangesKey).(*bool)
				}
				db.CloseTx(nil, changesMade)
				db.Tx = nil
				db.TxId = ""
			}
		}()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
