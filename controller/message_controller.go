package controller

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"go.chat/model"
	"go.chat/utils"
)

func AddMessage(w http.ResponseWriter, r *http.Request) {
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
	openChat := app.GetOpenChat(author)
	if openChat == nil {
		log.Printf("--%s-> AddMessage WARN no open chat for %s\n", utils.GetReqId(r), author)
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> AddMessage TRACE storing message for [%s] in [%s]\n", utils.GetReqId(r), author, openChat.Log())
	message, err := openChat.AddMessage(author, model.Message{ID: 0, Author: author, Text: msg})
	if err != nil {
		log.Printf("--%s-> AddMessage ERROR add message, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO distribute to all users in chat
	chatUsers, err := openChat.GetUsers(author)
	if err != nil || chatUsers == nil {
		log.Printf("--%s-> AddMessage ERROR get users, chat[%s], %s\n",
			utils.GetReqId(r), openChat.Log(), err)
	} else if len(chatUsers) == 0 {
		log.Printf("--%s-> AddMessage ERROR chatUsers are empty, chat[%s], %s\n",
			utils.GetReqId(r), openChat.Log(), err)
	} else {
		for _, user := range chatUsers {
			conn := app.getConn(user)
			if conn == nil {
				log.Printf("--%s-> AddChat ERROR cannot distribute message[%s] to user[%s], %s\n",
					utils.GetReqId(r), message.Log(), user, err)
			} else {
				log.Printf("--%s-> AddChat TRACE distributing message[%s] to user[%s]\n",
					utils.GetReqId(r), message.Log(), author)
				conn.Channel <- model.UserUpdate{
					Type: model.MessageUpdate,
					Chat: openChat,
					Msg:  message,
					User: user,
				}
			}
		}
	}

	log.Printf("--%s-> AddMessage TRACE templating [%s]\n", utils.GetReqId(r), message.Log())
	html, err := message.ToTemplate(author).GetHTML()
	if err != nil {
		log.Printf("<-%s-- AddMessage ERROR html, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- AddMessage TRACE serving html\n", utils.GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> DeleteMessage\n", utils.GetReqId(r))
	author, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		log.Printf("--%s-> DeleteMessage ERROR request method\n", utils.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return
	}

	openChat := app.GetOpenChat(author)
	if openChat == nil {
		log.Printf("--%s-> DeleteMessage ERROR open template for [%s]\n", utils.GetReqId(r), author)
		return
	}
	openChat.RemoveMessage(author, id)

	w.WriteHeader(http.StatusFound)
	w.Write(make([]byte, 0))
}
