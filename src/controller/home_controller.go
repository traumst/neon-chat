package controller

import (
	"log"
	"net/http"

	"prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/state"
	"prplchat/src/model/app"
	"prplchat/src/model/template"
	h "prplchat/src/utils/http"
)

func RenderLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("[%s] homeLogin TRACE IN", h.GetReqId(r))
	login := template.AuthTemplate{}
	home := template.HomeTemplate{
		Chats:         nil,
		OpenChat:      nil,
		User:          template.UserTemplate{UserName: "anon"},
		IsAuthorized:  false,
		LoginTemplate: login,
		Avatar:        nil,
	}
	log.Printf("[%s] homeLogin TRACE templating", h.GetReqId(r))
	html, err := home.HTML()
	if err != nil {
		log.Printf("[%s] homeLogin ERROR login %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home login"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func RenderHome(
	app *state.State,
	db *db.DBConn,
	w http.ResponseWriter,
	r *http.Request,
	user *app.User,
) {
	if app == nil {
		panic("app is nil")
	} else if db == nil {
		panic("db is nil")
	} else if user == nil {
		panic("user is nil")
	}
	log.Printf("[%s] homePage TRACE IN", h.GetReqId(r))
	var openChatTemplate *template.ChatTemplate
	var openChatId int = -1
	var openChatOwnerId uint = 0
	openChat := app.GetOpenChat(user.Id)
	if openChat == nil {
		log.Printf("[%s] homePage DEBUG, user[%d] has no open chat\n", h.GetReqId(r), user.Id)
		openChatTemplate = nil
	} else {
		log.Printf("[%s] homePage DEBUG, user[%d] has chat[%d] open\n", h.GetReqId(r), user.Id, openChat.Id)
		openChatTemplate = openChat.Template(user, user)
		openChatId = openChat.Id
		openChatOwnerId = openChat.Owner.Id
	}
	var chatTemplates []*template.ChatTemplate
	for _, chat := range app.GetChats(user.Id) {
		chatTemplates = append(chatTemplates, chat.Template(user, user))
	}
	var avatarTmpl *template.AvatarTemplate
	if dbAvatar, err := db.GetAvatar(user.Id); dbAvatar != nil && err == nil {
		avatar := handler.AvatarFromDB(*dbAvatar)
		avatarTmpl = avatar.Template(user)
	}
	home := template.HomeTemplate{
		Chats:         chatTemplates,
		OpenChat:      openChatTemplate,
		User:          *user.Template(openChatId, openChatOwnerId, user.Id),
		IsAuthorized:  true,
		LoginTemplate: template.AuthTemplate{},
		Avatar:        avatarTmpl,
	}
	html, err := home.HTML()
	if err != nil {
		log.Printf("[%s] homePage ERROR, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home page"))
		return
	}
	log.Printf("[%s] homePage TRACE, user[%d] gets content\n", h.GetReqId(r), user.Id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
