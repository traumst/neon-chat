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
	chat := app.State.GetOpenChat(author)
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

	distributeBetween(chat, author, html, r)

	log.Printf("<-%s-- AddMessage TRACE serving html\n", utils.GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func distributeBetween(chat *model.Chat, author string, html string, r *http.Request) {
	users, err := chat.GetUsers(author)
	if err != nil || users == nil {
		log.Printf("--%s-> distributeBetween ERROR get users, chat[%+v], %s\n",
			utils.GetReqId(r), chat, err)
		return
	}
	if len(users) == 0 {
		log.Printf("--%s-> distributeBetween ERROR chatUsers are empty, chat[%+v], %s\n",
			utils.GetReqId(r), chat, err)
		return
	}

	for _, user := range users {
		if user == author {
			log.Printf("--%s-> distributeBetween INFO new message is not sent to author[%s]\n",
				utils.GetReqId(r), user)
			continue
		}

		conn, err := app.State.GetConn(user)
		if err != nil {
			log.Printf("--%s-> distributeBetween ERROR cannot distribute html[%s] to user[%s], %s\n",
				utils.GetReqId(r), html, user, err)
			continue
		}

		log.Printf("--%s-> distributeBetween TRACE distributing html[%s] to user[%s]\n",
			utils.GetReqId(r), html, author)
		conn.Channel <- model.UserUpdate{
			Type:   model.MessageUpdate,
			ChatID: chat.ID,
			Author: author,
			Msg:    html,
		}
	}
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> DeleteMessage\n", utils.GetReqId(r))
	author, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		log.Printf("<-%s-- DeleteMessage ERROR request method\n", utils.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR parse id, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	openChat := app.State.GetOpenChat(author)
	if openChat == nil {
		log.Printf("<-%s-- DeleteMessage ERROR open template for [%s]\n", utils.GetReqId(r), author)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = openChat.DropMessage(author, id)
	if err != nil {
		log.Printf("<-%s-- DeleteMessage ERROR remove message[%d] from [%s], %s\n",
			utils.GetReqId(r), id, openChat.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- DeleteMessage done\n", utils.GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write(make([]byte, 0))
}
