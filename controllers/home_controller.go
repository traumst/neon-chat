package controllers

import (
	"log"
	"net/http"

	"go.chat/models"
	"go.chat/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Home", reqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("--%s-> Home WARN user, %s\n", reqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	openChat, err := chats.OpenTemplate(user)
	if err == nil {
		log.Printf("--%s-> Home DEBUG, user[%s] has chat[%d] open\n", reqId(r), user, openChat.Chat.ID)
	} else {
		log.Printf("--%s-> Home DEBUG, user[%s] has no open chat\n", reqId(r), user)
	}

	home := models.Home{
		OpenTemplate: openChat,
		Chats:        chats.GetChats(user),
		ActiveUser:   user,
	}
	html, err := home.GetHTML()
	if err != nil {
		log.Printf("--%s-> Home ERROR, %s\n", reqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	log.Printf("--%s-> Home TRACE, user[%s] gets content\n", reqId(r), user)
	w.Write([]byte(html))
}
