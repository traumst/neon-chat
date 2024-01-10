package controllers

import (
	"html/template"
	"net/http"

	"go.chat/models"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	openChat := chats.GetOpenChat()
	data := models.ChatData{
		Users:    []string{usernameCookie.Value},
		Messages: openChat.GetMessages(),
	}

	tmpl := template.Must(template.ParseFiles("views/chat.html"))
	tmpl.Execute(w, data)
}
