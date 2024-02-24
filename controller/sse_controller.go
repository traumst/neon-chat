package controller

import (
	"log"
	"net/http"
	"sync"

	"go.chat/model"
	"go.chat/utils"
)

var app = model.App{}

func PollUpdates(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Printf("<-%s-- PollUpdates TRACE does not provide %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- PollUpdates ERROR auth user[%s], %s\n", utils.GetReqId(r), user, err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	log.Printf("---%s--> PollUpdates TRACE IN polling updates for [%s]\n", utils.GetReqId(r), user)
	conn := app.State.ReplaceConn(w, *r, user)
	if conn == nil {
		log.Printf("<--%s--- PollUpdates ERROR conn not be established for [%s]\n", utils.GetReqId(r), user)
		return
	}
	defer app.State.DropConn(conn, user)
	utils.SetSseHeaders(&conn.Writer)
	log.Printf("---%s--> PollUpdates TRACE sse initiated for [%s]\n", utils.GetReqId(r), user)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.PollUpdatesForUser(conn, user)
	}()
	wg.Wait()
	log.Printf("<-%s-- PollUpdates TRACE OUT user[%s]\n", utils.GetReqId(r), user)
}
