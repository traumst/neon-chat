package controllers

import (
	"html/template"
	"net/http"

	"go.chat/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	openChat := chats.GetOpenChat()
	data := models.PageData{
		OpenChat: openChat,
		Chats:    chats.GetChatsCollapsed(),
		Username: usernameCookie.Value,
	}

	tmpl := template.Must(template.ParseFiles("views/home.html"))
	tmpl.Execute(w, data)
}
