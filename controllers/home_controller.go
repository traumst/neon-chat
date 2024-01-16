package controllers

import (
	"log"
	"net/http"

	"go.chat/handlers"
	"go.chat/models"
)

func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Home", reqId(r))
	user, err := handlers.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	home := models.Home{
		OpenChat:   chats.GetOpenChat(),
		Chats:      chats.GetChats(),
		ActiveUser: user,
	}
	html, err := home.GetHTML()
	if err != nil {
		log.Printf("--%s-> Home ERROR, %s\n", reqId(r), err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func FavIcon(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> FavIcon", reqId(r))
	http.ServeFile(w, r, "icons/favicon.ico")
}
