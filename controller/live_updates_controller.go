package controller

import (
	"log"
	"net/http"

	"go.chat/utils"
)

func PollUpdates(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> PollChats TRACE IN\n", utils.GetReqId(r))
	if r.Method != "GET" {
		log.Printf("<-%s-- PollChats TRACE does not provide %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- PollChats ERROR auth, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	app.PollUpdates(w, *r, user)
	log.Printf("<-%s-- PollChats TRACE OUT\n", utils.GetReqId(r))
}
