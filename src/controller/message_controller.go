package controller

import (
	"log"
	"net/http"

	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"neon-chat/src/event"
	"neon-chat/src/handler/parse"
	"neon-chat/src/handler/pub"
	"neon-chat/src/sse"
	"neon-chat/src/state"
	h "neon-chat/src/utils/http"
)

func QuoteMessage(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] '%s' '%s'\n", reqId, r.Method, r.RequestURI)
	if r.Method != "GET" {
		log.Printf("TRACE [%s] method '%s' not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	args, err := parse.ParseQueryString(r)
	if err != nil {
		log.Printf("ERROR [%s] parsing arguments, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid arguments"))
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	quote, err := pub.GetQuote(state, dbConn, user, args.ChatId, args.MsgId)
	if err != nil {
		log.Printf("ERROR [%s] quoting message[%d] in chat[%d], %s\n",
			reqId, args.ChatId, args.MsgId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to quote message"))
		return
	}
	tmpl, err := quote.Template(user)
	if err != nil {
		log.Printf("ERROR [%s] generating template, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to generate template"))
		return
	}
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("ERROR [%s] generating html, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to generate html"))
		return
	}
	log.Printf("TRACE [%s] serving html\n", reqId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] '%s' '%s'\n", reqId, r.Method, r.RequestURI)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] method '%s' not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		log.Printf("ERROR [%s] AddMessage request method\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("WARN [%s] AddMessage bad argument - chatid\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid chat id"))
		return
	}
	inputText, err := parse.ReadFormValueString(r, "msg")
	if err != nil || len(inputText) < 1 {
		log.Printf("WARN [%s] AddMessage bad argument - msg\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message too short"))
		return
	}
	quoteId, _ := parse.ReadFormValueUint(r, "quoteid")
	author := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, appMsg, err := pub.AddMessage(state, dbConn, chatId, author, inputText, quoteId)
	if err != nil || appMsg == nil {
		log.Printf("ERROR [%s] AddMessage while handing, %s, %v\n", reqId, err.Error(), inputText)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed adding message"))
		return
	}
	err = sse.DistributeMsg(state, dbConn.Tx, appChat, appMsg, event.MessageAdd)
	if err != nil {
		log.Printf("ERROR HandleMessageAdd distributing msg update, %s\n", err)
	}
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("TRACE [%s] AddMessage serving html\n", reqId)
	w.WriteHeader(http.StatusAccepted)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] '%s' '%s'\n", reqId, r.Method, r.RequestURI)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] method '%s' not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("ERROR [%s] bad arg - chatid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgId, err := parse.ReadFormValueUint(r, "msgid")
	if err != nil {
		log.Printf("ERROR [%s] bad arg - msgid, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, deletedMsg, err := pub.DeleteMessage(state, dbConn, chatId, user, msgId)
	if err != nil {
		log.Printf("ERROR [%s] deletion failed: %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if deletedMsg == nil {
		log.Printf("WARN [%s] message[%d] not found in chat[%d] for user[%d]\n",
			reqId, msgId, chatId, user.Id)
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}
	err = sse.DistributeMsg(state, dbConn.Tx, appChat, deletedMsg, event.MessageDrop)
	if err != nil {
		log.Printf("ERROR  distributing msg update, %s\n", err)
	}
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("TRACE [%s] DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("~~deleted~~"))
}
