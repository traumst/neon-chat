package controller

import (
	"log"
	"net/http"
	"sync"

	"neon-chat/src/consts"
	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	h "neon-chat/src/utils/http"
)

func PollUpdates(s *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	if r.Method != http.MethodGet {
		log.Printf("[%s] PollUpdates TRACE does not accept %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(s, db, w, r)
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

	h.SetSseHeaders(&conn.Writer)
	log.Printf("[%s] PollUpdates TRACE sse initiated for user[%d]\n", reqId, user.Id)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(s *state.State, conn *state.Conn) {
		defer s.DropConn(conn)
		defer wg.Done()
		handler.PollLiveUpdates(s, conn, user.Id)
		log.Printf("[%s] PollUpdates TRACE OUT user[%d]\n", reqId, user.Id)
	}(s, conn)
	wg.Wait()
}
