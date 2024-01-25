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
	log.Printf("--%s-> AddMessage\n", GetReqId(r))
	author, err := utils.GetCurrentUser(r)
	if err != nil || author == "" {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		log.Printf("--%s-> AddMessage ERROR request method\n", GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := template.HTMLEscapeString(r.FormValue("msg"))
	if msg == "" {
		log.Printf("--%s-> AddMessage WARN \n", GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	openChat, err := chats.OpenTemplate(author)
	if err != nil {
		log.Printf("--%s-> AddMessage ERROR open template, %s\n", GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if openChat == nil {
		log.Printf("--%s-> AddMessage ERROR openChat\n", GetReqId(r))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message, err := openChat.Chat.AddMessage(author, models.Message{ID: 0, Author: author, Text: msg})
	if err != nil {
		log.Printf("--%s-> AddMessage ERROR add message, %s\n", GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	html, err := message.GetHTML()
	if err != nil {
		log.Printf("<-%s-- AddMessage ERROR html, %s\n", GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- AddMessage TRACE serving html\n", GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> DeleteMessage\n", GetReqId(r))
	author, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		log.Printf("--%s-> DeleteMessage ERROR request method\n", GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := getMessageId(r)
	if err != nil {
		return
	}
	openChat, err := chats.OpenTemplate(author)
	if err != nil {
		log.Printf("--%s-> DeleteMessage ERROR open template for [%s]\n", GetReqId(r), author)
		return
	}
	openChat.Chat.RemoveMessage(author, id)

	w.WriteHeader(http.StatusFound)
	w.Write(make([]byte, 0))
}

func getMessageId(r *http.Request) (int, error) {
	return strconv.Atoi(r.FormValue("id"))
}
