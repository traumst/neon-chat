package controller

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"go.chat/handler"
	"go.chat/model"
	a "go.chat/model/app"
	e "go.chat/model/event"
	"go.chat/utils"
)

func AddMessage(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddMessage TRACE\n", utils.GetReqId(r))
	author, err := utils.GetCurrentUser(r)
	if err != nil || author == "" {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		log.Printf("--%s-> AddMessage ERROR request method\n", utils.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> AddMessage TRACE parsing input\n", utils.GetReqId(r))
	msg := template.HTMLEscapeString(r.FormValue("msg"))
	if msg == "" {
		log.Printf("--%s-> AddMessage WARN \n", utils.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> AddMessage TRACE opening current chat for [%s]\n", utils.GetReqId(r), author)
	// TODO verify open chat is modified chat
	chat := app.GetOpenChat(author)
	if chat == nil {
		log.Printf("--%s-> AddMessage WARN no open chat for %s\n", utils.GetReqId(r), author)
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> AddMessage TRACE storing message for [%s] in [%s]\n", utils.GetReqId(r), author, chat.Name)
	message, err := chat.AddMessage(author, a.Message{
		ID:     0,
		ChatID: chat.ID,
		Owner:  chat.Owner,
		Author: author,
		Text:   msg,
	})
	if err != nil {
		log.Printf("--%s-> AddMessage ERROR add message, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("--%s-> AddMessage TRACE templating message\n", utils.GetReqId(r))
	html, err := message.Template(author).HTML()
	if err != nil {
		log.Printf("<-%s-- AddMessage ERROR html [%+v], %s\n", utils.GetReqId(r), message, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.DistributeMsg(app, chat, author, message, e.MessageAdded)
	if err != nil {
		log.Printf("<-%s-- AddMessage ERROR distribute message, %s\n", utils.GetReqId(r), err)
	}

	log.Printf("<-%s-- AddMessage TRACE serving html\n", utils.GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteMessage(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteMessage\n", reqId)
	author, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	inChatID := r.PostFormValue("chatid")
	if inChatID == "" {
		log.Printf("<-%s-- DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatID, err := strconv.Atoi(inChatID)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR parse chatid, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat := app.GetOpenChat(author)
	if chat == nil {
		log.Printf("<-%s-- DeleteMessage ERROR open template for [%s]\n", reqId, author)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chatID != chat.ID {
		log.Printf("<-%s-- DeleteMessage ERROR chat id mismatch, %d != %d\n", reqId, chatID, chat.ID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inMsgID := r.PostFormValue("msgid")
	if inMsgID == "" {
		log.Printf("<-%s-- DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgID, err := strconv.Atoi(inMsgID)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR parse msgid, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, err := chat.DropMessage(author, msgID)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR remove message[%d] from [%s], %s\n", reqId, msgID, chat.Name, err)
		// TODO not necessarily StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO verify message delete distribution
	handler.DistributeMsg(app, chat, author, msg, e.MessageDeleted)

	log.Printf("<-%s-- DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("[deleted]"))
}
