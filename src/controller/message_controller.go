package controller

import (
	"log"
	"net/http"

	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	h "neon-chat/src/utils/http"
)

func QuoteMessage(
	state *state.State,
	db *d.DBConn,
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("[%s] QuoteMessage TRACE\n", h.GetReqId(r))
	if r.Method != "GET" {
		log.Printf("[%s] QuoteMessage ERROR request method\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	user, _ := handler.ReadSession(state, db, w, r)
	if user == nil {
		log.Printf("[%s] QuoteMessage ERROR user has no session\n", h.GetReqId(r))
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	args, err := handler.QueryStringArgs(r)
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR parsing arguments, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid arguments"))
		return
	}
	tmplMsg, err := handler.HandleGetMessage(state, db, user, args.ChatId, args.MsgId)
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR quoting message, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to quote message"))
		return
	}
	html, err := tmplMsg.ShortHTML()
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR templating quote, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to template message"))
		return
	}
	log.Printf("[%s] QuoteMessage TRACE serving html\n", h.GetReqId(r))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func AddMessage(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] TRACE AddMessage \n", h.GetReqId(r))
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
		log.Printf("[%s] WARN AddMessage bad argument - chatid\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid chat id"))
		return
	}
	msg, err := handler.FormValueString(r, "msg")
	if err != nil || len(msg) < 1 {
		log.Printf("[%s] WARN AddMessage bad argument - msg\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message too short"))
		return
	}
	quoteId, _ := handler.FormValueUint(r, "quoteid")
	message, err := handler.HandleMessageAdd(state, db, chatId, author, msg, quoteId)
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
