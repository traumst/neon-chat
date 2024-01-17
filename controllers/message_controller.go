package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"go.chat/models"
	"go.chat/utils"
)

func AddMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddMessage\n", reqId(r))
	author, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		log.Printf("--%s-> AddMessage ERROR request method\n", reqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	text := template.HTMLEscapeString(r.FormValue("text"))
	if text == "" || author == "" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}

	openChat, err := chats.OpenTemplate(author)
	if openChat == nil {
		log.Printf("--%s-> AddMessage ERROR openChat\n", reqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err := openChat.Chat.AddMessage(author, models.Message{ID: 0, Author: author, Text: text})
	if err != nil {
		log.Printf("--%s-> AddMessage ERROR add message, %s\n", reqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	html, err := msg.GetHTML()
	if err != nil {
		log.Printf("<-%s-- AddMessage ERROR html, %s\n", reqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(html))
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> DeleteMessage\n", reqId(r))
	author, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		log.Printf("--%s-> DeleteMessage ERROR request method\n", reqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := getMessageId(r)
	if err != nil {
		return
	}
	openChat, err := chats.OpenTemplate(author)
	if openChat == nil {
		return
	}
	openChat.Chat.RemoveMessage(author, id)

	w.Write(make([]byte, 0))
}

func getMessageId(r *http.Request) (int, error) {
	return strconv.Atoi(r.FormValue("id"))
}
