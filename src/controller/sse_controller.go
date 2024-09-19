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

func PollUpdates(s *state.State, dbConn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	if r.Method != http.MethodGet {
		log.Printf("[%s] PollUpdates TRACE does not accept %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := pub.ReadSession(s, dbConn, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] PollUpdates WARN user, %s\n", reqId, err)
		return
	}
	log.Printf("[%s] PollUpdates TRACE IN polling updates for user[%d]\n", reqId, user.Id)
	conn := s.AddConn(w, *r, user, nil)
	if conn == nil {
		log.Printf("[%s] PollUpdates ERROR conn not be established for user[%d]\n", reqId, user.Id)
		return
	}
	defer s.DropConn(conn)

	h.SetSseHeaders(&conn.Writer)
	log.Printf("[%s] PollUpdates TRACE sse initiated for user[%d]\n", reqId, user.Id)

	sse.LiveUpdates(s, conn, user.Id)
	log.Printf("[%s] PollUpdates TRACE OUT user[%d]\n", reqId, user.Id)
}
