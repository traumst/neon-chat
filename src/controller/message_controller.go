package controller

import (
	"log"
	"net/http"
	"strconv"

	"go.chat/src/handler"
	a "go.chat/src/model/app"
	"go.chat/src/model/event"
	"go.chat/src/utils"
)

func AddMessage(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddMessage TRACE\n", utils.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("--%s-> AddMessage ERROR request method\n", utils.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	author, err := handler.ReadSession(app, w, r)
	if err != nil || author == nil {
		RenderHome(app, w, r)
		return
	}
	log.Printf("--%s-> AddMessage TRACE parsing input\n", utils.GetReqId(r))
	msg := r.FormValue("msg")
	msg = utils.TrimSpaces(msg)
	msg = utils.TrimSpecial(msg)
	if msg == "" {
		log.Printf("--%s-> AddMessage WARN \n", utils.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> AddMessage TRACE opening current chat for user[%d]\n", utils.GetReqId(r), author.Id)
	// TODO verify open chat is modified chat
	chat := app.GetOpenChat(author.Id)
	if chat == nil {
		log.Printf("--%s-> AddMessage WARN no open chat for user[%d]\n", utils.GetReqId(r), author.Id)
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> AddMessage TRACE storing message for user[%d] in chat[%s]\n",
		utils.GetReqId(r), author.Id, chat.Name)
	message, err := chat.AddMessage(author.Id, a.Message{
		Id:     0,
		ChatId: chat.Id,
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

	err = handler.DistributeMsg(app, chat, author.Id, message, event.MessageAdded)
	if err != nil {
		log.Printf("<-%s-- AddMessage ERROR distribute message, %s\n", utils.GetReqId(r), err)
	}

	log.Printf("<-%s-- AddMessage TRACE serving html\n", utils.GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteMessage(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteMessage\n", reqId)
	author, err := handler.ReadSession(app, w, r)
	if err != nil || author == nil {
		RenderHome(app, w, r)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	inChatId := r.PostFormValue("chatid")
	if inChatId == "" {
		log.Printf("<-%s-- DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId, err := strconv.Atoi(inChatId)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR parse chatid, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat := app.GetOpenChat(author.Id)
	if chat == nil {
		log.Printf("<-%s-- DeleteMessage ERROR open template for user[%d]\n", reqId, author.Id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chatId != chat.Id {
		log.Printf("<-%s-- DeleteMessage ERROR chat id mismatch, %d != %d\n", reqId, chatId, chat.Id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inMsgId := r.PostFormValue("msgid")
	if inMsgId == "" {
		log.Printf("<-%s-- DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msgId, err := strconv.Atoi(inMsgId)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR parse msgid, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, err := chat.DropMessage(author.Id, msgId)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR remove message[%d] from [%s], %s\n", reqId, msgId, chat.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = handler.DistributeMsg(app, chat, author.Id, msg, event.MessageDeleted)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR distribute message, %s\n", reqId, err)
	}

	log.Printf("<-%s-- DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("[DeletedM]"))
}
