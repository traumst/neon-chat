package controllers

import (
	"html/template"
	"net/http"

	"go.chat/models"
)

func OpenChat(w http.ResponseWriter, r *http.Request) {
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

func AddChat(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	chatName := r.FormValue("chatName")
	chats.AddChat(usernameCookie.Value, chatName)

	http.Redirect(w, r, "/", http.StatusFound)
}
