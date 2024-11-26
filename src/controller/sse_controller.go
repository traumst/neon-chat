package controller

import (
	"log"
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/db"
	"neon-chat/src/handler/pub"
	"neon-chat/src/sse"
	"neon-chat/src/state"
	h "neon-chat/src/utils/http"
)

func PollUpdates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqId := ctx.Value(consts.ReqIdKey).(string)
	if r.Method != http.MethodGet {
		log.Printf("TRACE [%s] '%s' does not accept %s\n", reqId, r.RequestURI, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s := ctx.Value(consts.AppState).(*state.State)
	dbConn := ctx.Value(consts.DBConn).(*db.DBConn)

	user, err := pub.ReadSession(s, dbConn, w, r)
	if err != nil || user == nil {
		log.Printf("WARN [%s]  user, %s\n", reqId, err)
		return
	}
	// WIP consider instead
	//dbConn := ctx.Value(consts.DBConn).(*db.DBConn)
	//user, err := pub.ReadSession(s, dbConn, w, r)
	// user := ctx.Value(consts.ActiveUser).(*app.User)
	// if user == nil {
	// 	log.Printf("WARN [%s] user is nil\n", reqId)
	// 	return
	// }
	log.Printf("TRACE [%s] polling updates for user[%d]\n", reqId, user.Id)
	conn := s.AddConn(w, *r, user, nil)
	if conn == nil {
		log.Printf("WARN [%s] conn not be established for user[%d]\n", reqId, user.Id)
		return
	}
	defer s.DropConn(conn)

	h.SetSseHeaders(&conn.Writer)
	log.Printf("TRACE [%s] sse initiated for user[%d]\n", reqId, user.Id)

	sse.PollUpdates(s, conn, user.Id)
	log.Printf("TRACE [%s] live update consumption stopped for user[%d]\n", reqId, user.Id)
}
