package controller

import (
	"log"
	"net/http"

	"neon-chat/src/consts"
	d "neon-chat/src/db"
	"neon-chat/src/handler"
	pi "neon-chat/src/handler/parse"
	a "neon-chat/src/model/app"
	"neon-chat/src/state"
	h "neon-chat/src/utils/http"
)

func QuoteMessage(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] QuoteMessage TRACE\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] QuoteMessage ERROR request method\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	args, err := pi.ParseQueryString(r)
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR parsing arguments, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid arguments"))
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*a.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	db := r.Context().Value(consts.DBConn).(*d.DBConn)
	html, err := handler.HandleGetQuote(state, db, user, args.ChatId, args.MsgId)
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR quoting message[%d] in chat[%d], %s\n",
			reqId, args.ChatId, args.MsgId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to quote message"))
		return
	}
	log.Printf("[%s] QuoteMessage TRACE serving html\n", reqId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] TRACE AddMessage \n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] AddMessage ERROR request method\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	chatId, err := pi.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] WARN AddMessage bad argument - chatid\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid chat id"))
		return
	}
	msg, err := pi.ReadFormValueString(r, "msg")
	if err != nil || len(msg) < 1 {
		log.Printf("[%s] WARN AddMessage bad argument - msg\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message too short"))
		return
	}
	quoteId, _ := pi.ReadFormValueUint(r, "quoteid")
	author := r.Context().Value(consts.ActiveUser).(*a.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	db := r.Context().Value(consts.DBConn).(*d.DBConn)
	message, err := handler.HandleMessageAdd(state, db, chatId, author, msg, quoteId)
	if err != nil || message == nil {
		log.Printf("[%s] AddMessage ERROR while handing, %s, %v\n", reqId, err.Error(), message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed adding message"))
		return
	}
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("[%s] AddMessage TRACE serving html\n", reqId)
	w.WriteHeader(http.StatusAccepted)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] DeleteMessage\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatId, err := pi.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - chatid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgId, err := pi.ReadFormValueUint(r, "msgid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - msgid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*a.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	db := r.Context().Value(consts.DBConn).(*d.DBConn)
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
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("[%s] DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("~~deleted~~"))
}
