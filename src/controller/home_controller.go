package controller

import (
	"log"
	"net/http"

	"go.chat/src/handler"
	"go.chat/src/model/event"
	"go.chat/src/model/template"
	"go.chat/src/utils"
)

func RenderHome(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Home", utils.GetReqId(r))
	user, err := handler.ReadSession(app, w, r)
	if err != nil {
		log.Printf("--%s-> Home INFO user is not authorized, %s\n", utils.GetReqId(r), err)
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
	openChatId := -1
	if openChatTemplate != nil {
		openChatId = openChat.Id
	}
	home := template.HomeTemplate{
		OpenTemplate: openChatTemplate,
		Chats:        chatTemplates,
		ActiveUser:   user.Name,
		LoadLocal:    app.LoadLocal(),
		ChatAddEvent: event.ChatAddEventName.Format(openChatId, user.Id, -5),
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
