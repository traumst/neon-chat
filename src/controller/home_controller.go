package controller

import (
	"log"
	"net/http"

	"go.chat/src/handler"
	"go.chat/src/model/app"
	"go.chat/src/model/event"
	"go.chat/src/model/template"
	h "go.chat/src/utils/http"
)

func RenderHome(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> RenderHome", h.GetReqId(r))
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
		Chats:          nil,
		OpenTemplate:   nil,
		ActiveUser:     "User",
		LoadLocal:      app.LoadLocal(),
		IsAuthorized:   false,
		LoginTemplate:  login,
		ChatAddEvent:   "",
		ChatCloseEvent: "",
	}
	html, err := home.HTML()
	if err != nil {
		log.Printf("--%s-> homeLogin ERROR login %s\n", h.GetReqId(r), err)
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
		log.Printf("--%s-> homePage DEBUG, user[%d] has no open chat\n", h.GetReqId(r), user.Id)
		openChatTemplate = nil
	} else {
		log.Printf("--%s-> homePage DEBUG, user[%d] has chat[%d] open\n", h.GetReqId(r), user.Id, openChat.Id)
		openChatTemplate = openChat.Template(user, user)
	}
	var chatTemplates []*template.ChatTemplate
	for _, chat := range app.GetChats(user.Id) {
		chatTemplates = append(chatTemplates, chat.Template(user, user))
	}
	openChatId := -1
	if openChatTemplate != nil {
		openChatId = openChat.Id
	}
	home := template.HomeTemplate{
		Chats:          chatTemplates,
		OpenTemplate:   openChatTemplate,
		ActiveUser:     user.Name,
		LoadLocal:      app.LoadLocal(),
		IsAuthorized:   true,
		ChatAddEvent:   event.ChatAdd.FormatEventName(openChatId, user.Id, -5),
		ChatCloseEvent: event.ChatClose.FormatEventName(openChatId, user.Id, -6),
	}
	html, err := home.HTML()
	if err != nil {
		log.Printf("--%s-> homePage ERROR, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home page"))
		return
	}
	log.Printf("--%s-> homePage TRACE, user[%d] gets content\n", h.GetReqId(r), user.Id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
