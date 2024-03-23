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
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> Home WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	var openChatTemplate *template.ChatTemplate
	openChat := app.GetOpenChat(user.Id)
	if openChat == nil {
		log.Printf("--%s-> Home DEBUG, user[%d] has no open chat\n", utils.GetReqId(r), user.Id)
		openChatTemplate = nil
	} else {
		log.Printf("--%s-> Home DEBUG, user[%d] has chat[%d] open\n", utils.GetReqId(r), user.Id, openChat.Id)
		openChatTemplate = openChat.Template(user)
	}

	var chatTemplates []*template.ChatTemplate
	for _, chat := range app.GetChats(user.Id) {
		chatTemplates = append(chatTemplates, chat.Template(user))
	}

	home := template.HomeTemplate{
		OpenTemplate: openChatTemplate,
		Chats:        chatTemplates,
		ActiveUser:   user.Name,
		LoadLocal:    app.LoadLocal(),
	}
	html, err := home.HTML()
	if err != nil {
		log.Printf("--%s-> Home ERROR, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	log.Printf("--%s-> Home TRACE, user[%d] gets content\n", utils.GetReqId(r), user.Id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
