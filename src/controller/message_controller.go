package controller

import (
	"log"
	"net/http"

	d "prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/state"
	h "prplchat/src/utils/http"
)

func AddMessage(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] AddMessage TRACE\n", h.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("[%s] AddMessage ERROR request method\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	author, err := handler.ReadSession(state, db, w, r)
	if err != nil || author == nil {
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] AddMessage WARN bad argument - chatid\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid chat id"))
		return
	}
	msg, err := handler.FormValueString(r, "msg")
	if err != nil || len(msg) < 1 {
		log.Printf("[%s] AddMessage WARN bad argument - msg\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message too short"))
		return
	}

	message, err := handler.HandleMessageAdd(state, db, chatId, author, msg)
	if err != nil || message == nil {
		log.Printf("[%s] AddMessage ERROR while handing, %s, %v\n", h.GetReqId(r), err.Error(), message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed adding message"))
		return
	}

	log.Printf("[%s] AddMessage TRACE serving html\n", h.GetReqId(r))
	w.WriteHeader(http.StatusAccepted)
}

func DeleteMessage(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] DeleteMessage\n", reqId)
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - chatid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgId, err := handler.FormValueUint(r, "msgid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - msgid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deleted, err := handler.HandleMessageDelete(state, db, chatId, user, msgId)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR deletion failed: %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deleted == nil {
		log.Printf("[%s] DeleteMessage WARN message[%d] not found in chat[%d] for user[%d]\n",
			reqId, msgId, chatId, user.Id)
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}

	log.Printf("[%s] DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("~~deleted~~"))
}
