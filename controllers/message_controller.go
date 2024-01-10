package controllers

import (
	"html/template"
	"net/http"
	"strconv"

	"go.chat/models"
)

func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		author := usernameCookie.Value
		text := template.HTMLEscapeString(r.FormValue("text"))
		if text != "" && author != "" {
			openChat := chats.GetOpenChat()
			if openChat != nil {
				openChat.AddMessage(models.Message{Author: author, Text: text})
			}
		}
	}
	// redirect to re-render chat history
	http.Redirect(w, r, "/", http.StatusFound)
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, err := strconv.Atoi(r.FormValue("id"))
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
