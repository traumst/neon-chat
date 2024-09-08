package controller

import (
	"log"
	"net/http"
	"sync"

	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	h "neon-chat/src/utils/http"
)

func PollUpdates(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Printf("[%s] PollUpdates TRACE does not provide %s\n", h.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] PollUpdates WARN user, %s\n", h.GetReqId(r), err)
		return
	}
	log.Printf("[%s] PollUpdates TRACE IN polling updates for user[%d]\n", h.GetReqId(r), user.Id)
	conn := state.AddConn(w, *r, user, nil)
	if conn == nil {
		log.Printf("[%s] PollUpdates ERROR conn not be established for user[%d]\n", h.GetReqId(r), user.Id)
		return
	}
	defer state.DropConn(conn)
	h.SetSseHeaders(&conn.Writer)
	log.Printf("[%s] PollUpdates TRACE sse initiated for user[%d]\n", h.GetReqId(r), user.Id)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handler.PollLiveUpdates(state, conn, user.Id)
	}()
	wg.Wait()
	log.Printf("[%s] PollUpdates TRACE OUT user[%d]\n", h.GetReqId(r), user.Id)
}
