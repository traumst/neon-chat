package controller

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"go.chat/handler"
	"go.chat/model"
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
	chat := app.GetOpenChat(author)
	if chat == nil {
		log.Printf("--%s-> AddMessage WARN no open chat for %s\n", utils.GetReqId(r), author)
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> AddMessage TRACE storing message for [%s] in [%s]\n", utils.GetReqId(r), author, chat.Name)
	message, err := chat.AddMessage(author, model.Message{ID: 0, Author: author, Text: msg})
	if err != nil {
		log.Printf("--%s-> AddMessage ERROR add message, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("--%s-> AddMessage TRACE templating message\n", utils.GetReqId(r))
	html, err := message.ToTemplate(author).GetHTML()
	if err != nil {
		log.Printf("<-%s-- AddMessage ERROR html [%+v], %s\n", utils.GetReqId(r), message, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.DistributeMsg(app, chat, author, message, model.MessageAdded)

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
	chatID := r.PostFormValue("msgid")
	if chatID == "" {
		log.Printf("<-%s-- DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chatID)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR parse id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat := app.GetOpenChat(author)
	if chat == nil {
		log.Printf("<-%s-- DeleteMessage ERROR open template for [%s]\n", reqId, author)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg, err := chat.DropMessage(author, id)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR remove message[%d] from [%s], %s\n", reqId, id, chat.Name, err)
		// TODO not necessarily StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO verify message delete distribution
	handler.DistributeMsg(app, chat, author, msg, model.MessageDeleted)

	log.Printf("<-%s-- DeleteMessage done\n", reqId)
	w.WriteHeader(http.StatusAccepted)
	w.Write(make([]byte, 0))
}
