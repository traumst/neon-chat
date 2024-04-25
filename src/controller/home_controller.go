package controller

import (
	"log"
	"net/http"

	"go.chat/src/handler"
	"go.chat/src/model/app"
	"go.chat/src/model/event"
	"go.chat/src/model/template"
	"go.chat/src/utils"
)

func RenderHome(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> RenderHome", utils.GetReqId(r))
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		homeLogin(app, w, r)
	} else {
		homePage(app, w, r, user)
	}
}

func homeLogin(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	login := template.LoginTemplate{
		Login: template.AuthForm{
			Id:    "login",
			Label: "Login",
			Title: "Login",
		},
		Signup: template.AuthForm{
			Id:    "signup",
			Label: "Signup",
			Title: "Signup",
		},
		LoadLocal: app.LoadLocal(),
	}
	home := template.HomeTemplate{
		OpenTemplate:  nil,
		Chats:         nil,
		ActiveUser:    "",
		LoadLocal:     app.LoadLocal(),
		ChatAddEvent:  "",
		IsAuthorized:  false,
		LoginTemplate: login,
	}
	html, err := home.HTML()
	if err != nil {
		log.Printf("--%s-> homeLogin ERROR login %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home login"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func homePage(app *handler.AppState, w http.ResponseWriter, r *http.Request, user *app.User) {
	var openChatTemplate *template.ChatTemplate
	openChat := app.GetOpenChat(user.Id)
	if openChat == nil {
		log.Printf("--%s-> homePage DEBUG, user[%d] has no open chat\n", utils.GetReqId(r), user.Id)
		openChatTemplate = nil
	} else {
		log.Printf("--%s-> homePage DEBUG, user[%d] has chat[%d] open\n", utils.GetReqId(r), user.Id, openChat.Id)
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
		IsAuthorized: true,
	}
	html, err := home.HTML()
	if err != nil {
		log.Printf("--%s-> homePage ERROR, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home page"))
		return
	}
	log.Printf("--%s-> homePage TRACE, user[%d] gets content\n", utils.GetReqId(r), user.Id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
