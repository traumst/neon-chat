package controller

import (
	"log"
	"net/http"
	"sync"

	"go.chat/src/handler"
	"go.chat/src/utils"
)

func PollUpdates(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Printf("<-%s-- PollUpdates TRACE does not provide %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("--%s-> PollUpdates WARN user, %s\n", utils.GetReqId(r), err)
		return
	}
	log.Printf("---%s--> PollUpdates TRACE IN polling updates for user[%d]\n", utils.GetReqId(r), user.Id)
	conn := app.ReplaceConn(w, *r, user)
	if conn == nil {
		log.Printf("<--%s--- PollUpdates ERROR conn not be established for user[%d]\n", utils.GetReqId(r), user.Id)
		return
	}
	defer app.DropConn(conn)
	utils.SetSseHeaders(&conn.Writer)
	log.Printf("---%s--> PollUpdates TRACE sse initiated for user[%d]\n", utils.GetReqId(r), user.Id)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handler.PollUpdatesForUser(app, conn, user.Id)
	}()
	wg.Wait()
	log.Printf("<-%s-- PollUpdates TRACE OUT user[%d]\n", utils.GetReqId(r), user.Id)
}
