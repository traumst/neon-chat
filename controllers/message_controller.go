package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"go.chat/handlers"
	"go.chat/models"
)

func AddMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddMessage\n", reqId(r))
	author, err := handlers.GetCurrentUser(r)
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

	openChat := chats.GetOpenChat()
	if openChat == nil {
		log.Printf("--%s-> AddMessage ERROR openChat\n", reqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := openChat.AddMessage(models.Message{ID: 0, Author: author, Text: text})
	html, err := msg.GetHTML()
	if err != nil {
		log.Printf("<-%s-- GetMessage ERROR html, %s\n", reqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> GetMessage\n", reqId(r))
	_, err := handlers.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	msgId, err := getMessageId(r)
	if err != nil {
		log.Printf("--%s-> GetMessage ERROR id, %s\n", reqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := chats.GetOpenChat().GetMessage(msgId)
	html, err := msg.GetHTML()
	if err != nil {
		log.Printf("<-%s-- GetMessage ERROR html, %s\n", reqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> DeleteMessage\n", reqId(r))
	_, err := handlers.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {
		id, err := getMessageId(r)
		if err != nil {
			return
		}
		chat := chats.GetOpenChat()
		if chat == nil {
			return
		}
		chat.RemoveMessage(id)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func getMessageId(r *http.Request) (int, error) {
	return strconv.Atoi(r.FormValue("id"))
}
