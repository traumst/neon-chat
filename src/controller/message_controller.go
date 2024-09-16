package controller

import (
	"log"
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/parse"
	"neon-chat/src/model/app"
	"neon-chat/src/model/event"
	"neon-chat/src/sse"
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
	args, err := parse.ParseQueryString(r)
	if err != nil {
		log.Printf("[%s] QuoteMessage ERROR parsing arguments, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid arguments"))
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	html, err := handler.GetQuote(state, dbConn, user, args.ChatId, args.MsgId)
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
		log.Printf("[%s] ERROR AddMessage request method\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] WARN AddMessage bad argument - chatid\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid chat id"))
		return
	}
	inputText, err := parse.ReadFormValueString(r, "msg")
	if err != nil || len(inputText) < 1 {
		log.Printf("[%s] WARN AddMessage bad argument - msg\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message too short"))
		return
	}
	quoteId, _ := parse.ReadFormValueUint(r, "quoteid")
	author := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, appMsg, err := handler.AddMessage(state, dbConn, chatId, author, inputText, quoteId)
	if err != nil || appMsg == nil {
		log.Printf("[%s] ERROR AddMessage while handing, %s, %v\n", reqId, err.Error(), inputText)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed adding message"))
		return
	}
	err = sse.DistributeMsg(state, dbConn.Tx, appChat, appMsg, event.MessageAdd)
	if err != nil {
		log.Printf("ERROR HandleMessageAdd distributing msg update, %s\n", err)
	}
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("[%s] TRACE AddMessage serving html\n", reqId)
	w.WriteHeader(http.StatusAccepted)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] DeleteMessage\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - chatid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgId, err := parse.ReadFormValueUint(r, "msgid")
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR bad arg - msgid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, deletedMsg, err := handler.DeleteMessage(state, dbConn, chatId, user, msgId)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR deletion failed: %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deletedMsg == nil {
		log.Printf("[%s] DeleteMessage WARN message[%d] not found in chat[%d] for user[%d]\n",
			reqId, msgId, chatId, user.Id)
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}
	err = sse.DistributeMsg(state, dbConn.Tx, appChat, deletedMsg, event.MessageDrop)
	if err != nil {
		log.Printf("HandleMessageDelete ERROR distributing msg update, %s\n", err)
	}
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("[%s] DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("~~deleted~~"))
}
