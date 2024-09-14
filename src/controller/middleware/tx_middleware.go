package middleware

import (
	"context"
	"fmt"
	"log"
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
	"net/http"
)

func TransactionMiddleware() Middleware {
	return Middleware{
		Name: "Transaction",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				reqId := ctx.Value(consts.ReqIdKey).(string)
				// copy db connection
				db := *ctx.Value(consts.DBConn).(*db.DBConn)
				// populates tx prop
				_, txId, err := db.OpenTx(reqId)
				if err != nil {
					log.Printf("FATAL [%s] TransactionMiddleware failed to open transaction: %s", reqId, err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				} else if txId == "" || txId != reqId {
					log.Printf("WARN [%s] TransactionMiddleware unexpected txId: %s", reqId, txId)
				}
				// ctx with both general conn and per-session tx
				ctx = context.WithValue(ctx, consts.DBConn, &db)
				defer func() {
					if p := recover(); p != nil {
						log.Printf("FATAL [%s] TransactionMiddleware Failed to open transaction: %v", reqId, p)
						db.CloseTx(fmt.Errorf("panic: %v", p), false)
						panic(p) // re-throw the panic after rollback
					} else if code := w.(*h.StatefulWriter).Status(); code >= http.StatusBadRequest {
						db.CloseTx(fmt.Errorf("error status code %d", code), false)
					} else {
						db.CloseTx(nil, utils.HasTxChanges(r))
						db.Tx = nil
						db.TxId = ""
					}
				}()
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}}
}
