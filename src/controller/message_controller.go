package controller

import (
	"log"
	"net/http"
	"strconv"

	d "prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	"prplchat/src/model/event"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func AddMessage(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] AddMessage TRACE\n", h.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("[%s] AddMessage ERROR request method\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("action not allowed"))
		return
	}
	author, err := handler.ReadSession(app, db, w, r)
	if err != nil || author == nil {
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	log.Printf("[%s] AddMessage TRACE parsing input\n", h.GetReqId(r))
	addToChatId, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("[%s] AddMessage WARN \n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid chat id"))
		return
	}
	chatId := uint(addToChatId)
	msg := r.FormValue("msg")
	msg = utils.ReplaceWithSingleSpace(msg)
	if len(msg) < 1 {
		log.Printf("[%s] AddMessage WARN \n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message too short"))
		return
	}

	message, err := handler.HandleMessageAdd(app, db, chatId, author, msg)
	if err != nil || message == nil {
		log.Printf("[%s] AddMessage ERROR while handing, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed adding message"))
		return
	}

	log.Printf("[%s] AddMessage TRACE serving html\n", h.GetReqId(r))
	w.WriteHeader(http.StatusAccepted)
}

func DeleteMessage(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] DeleteMessage\n", reqId)
	author, err := handler.ReadSession(app, db, w, r)
	if err != nil || author == nil {
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	inChatId := r.PostFormValue("chatid")
	if inChatId == "" {
		log.Printf("[%s] DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currChatId, err := strconv.Atoi(inChatId)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR parse chatid, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId := uint(currChatId)
	chat := app.GetOpenChat(author.Id)
	if chat == nil {
		log.Printf("[%s] DeleteMessage ERROR open template for user[%d]\n", reqId, author.Id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chatId != chat.Id {
		log.Printf("[%s] DeleteMessage ERROR chat id mismatch, %d != %d\n", reqId, chatId, chat.Id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inMsgId := r.PostFormValue("msgid")
	if inMsgId == "" {
		log.Printf("[%s] DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currMsgId, err := strconv.Atoi(inMsgId)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR parse msgid, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgId := uint(currMsgId)
	msg, err := chat.DropMessage(author.Id, msgId)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR remove message[%d] from [%s], %s\n", reqId, msgId, chat.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = sse.DistributeMsg(app, chat, author.Id, msg, event.MessageDrop)
	if err != nil {
		log.Printf("[%s] DeleteMessage ERROR distribute message, %s\n", reqId, err)
	}

	log.Printf("[%s] DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusAccepted)
}
