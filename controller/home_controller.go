package controller

import (
	"log"
	"net/http"

	"go.chat/model"
	"go.chat/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Home", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("--%s-> Home WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	var template *model.ChatTemplate
	openChat := app.state.GetOpenChat(user)
	if openChat == nil {
		log.Printf("--%s-> Home DEBUG, user[%s] has no open chat\n", utils.GetReqId(r), user)
		template = nil
	} else {
		log.Printf("--%s-> Home DEBUG, user[%s] has chat[%d] open\n", utils.GetReqId(r), user, openChat.ID)
		template = openChat.ToTemplate(user)
	}

	home := model.HomeTemplate{
		OpenTemplate: template,
		Chats:        app.state.GetChats(user),
		ActiveUser:   user,
	}
	html, err := home.GetHTML(utils.GetReqId(r))
	if err != nil {
		log.Printf("--%s-> Home ERROR, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	log.Printf("--%s-> Home TRACE, user[%s] gets content\n", utils.GetReqId(r), user)
	w.Write([]byte(html))
}
