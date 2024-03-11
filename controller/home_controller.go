package controller

import (
	"log"
	"net/http"

	"go.chat/model"
	"go.chat/model/template"
	"go.chat/utils"
)

func Home(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Home", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("--%s-> Home WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	var openChatTemplate *template.ChatTemplate
	openChat := app.GetOpenChat(user)
	if openChat == nil {
		log.Printf("--%s-> Home DEBUG, user[%s] has no open chat\n", utils.GetReqId(r), user)
		openChatTemplate = nil
	} else {
		log.Printf("--%s-> Home DEBUG, user[%s] has chat[%d] open\n", utils.GetReqId(r), user, openChat.ID)
		openChatTemplate = openChat.ToTemplate(user)
	}

	var chatTemplates []*template.ChatTemplate
	for _, chat := range app.GetChats(user) {
		chatTemplates = append(chatTemplates, chat.ToTemplate(user))
	}

	home := template.HomeTemplate{
		OpenTemplate: openChatTemplate,
		Chats:        chatTemplates,
		ActiveUser:   user,
	}
	html, err := home.GetHTML()
	if err != nil {
		log.Printf("--%s-> Home ERROR, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	log.Printf("--%s-> Home TRACE, user[%s] gets content\n", utils.GetReqId(r), user)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
