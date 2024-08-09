package controller

import (
	"log"
	"net/http"

	"prplchat/src/convert"
	"prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/template"
	h "prplchat/src/utils/http"
)

func RenderLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Printf("[%s] RenderLogin TRACE IN", h.GetReqId(r))
	login := template.AuthTemplate{}
	home := template.HomeTemplate{
		Chats:         nil,
		OpenChat:      nil,
		User:          template.UserTemplate{UserName: "anon"},
		IsAuthorized:  false,
		LoginTemplate: login,
		Avatar:        nil,
	}
	log.Printf("[%s] RenderLogin TRACE templating", h.GetReqId(r))
	html, err := home.HTML()
	if err != nil {
		log.Printf("[%s] RenderLogin ERROR login %s\n", h.GetReqId(r), err)
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
	user *a.User,
) {
	if app == nil {
		panic("app is nil")
	} else if db == nil {
		panic("db is nil")
	} else if user == nil {
		panic("user is nil")
	}
	log.Printf("[%s] RenderHome TRACE IN", h.GetReqId(r))
	html, err := templateHome(app, db, r, user)
	if err != nil {
		log.Printf("[%s] RenderHome ERROR, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render home page"))
		return
	}
	log.Printf("[%s] RenderHome TRACE, user[%d] gets content\n", h.GetReqId(r), user.Id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func templateOpenChat(app *state.State, db *db.DBConn, user *a.User) *template.ChatTemplate {
	openChatId := app.GetOpenChat(user.Id)
	if openChatId == 0 {
		log.Printf("templateOpenchat DEBUG, user[%d] has no open chat\n", user.Id)
		return nil
	}
	openChat, err := db.GetChat(openChatId)
	if err != nil {
		log.Printf("templateOpenchat ERROR, failed to get chat[%d], %s\n", openChatId, err.Error())
		return nil // TODO custom error pop-up
	}
	dbChatUsers, err := db.GetChatUsers(user.Id)
	if err != nil {
		log.Printf("templateHome ERROR, failed getting chat[%d] users, %s\n", openChatId, err.Error())
		return nil
	}
	appChatUsers := make([]*a.User, 0)
	for _, dbUser := range dbChatUsers {
		appChatUsers = append(appChatUsers, convert.UserDBToApp(&dbUser))
	}
	appChat := convert.ChatDBToApp(openChat)
	return appChat.Template(user, user, appChatUsers)
}

func templateHome(
	app *state.State,
	db *db.DBConn,
	r *http.Request,
	user *a.User,
) (string, error) {
	var avatarTmpl *template.AvatarTemplate
	if dbAvatar, err := db.GetAvatar(user.Id); dbAvatar != nil && err == nil {
		avatar := convert.AvatarDBToApp(dbAvatar)
		avatarTmpl = avatar.Template(user)
	}

	openChatTemplate := templateOpenChat(app, db, user)
	chats, err := handler.GetChats(app, db, user.Id)
	if err != nil {
		log.Printf("[%s] templateHome ERROR, failed getting chats for user[%d], %s\n",
			h.GetReqId(r), user.Id, err.Error())
		return "", err
	}
	var chatTemplates []*template.ChatTemplate
	for _, chat := range chats {
		chatTemplates = append(chatTemplates, chat.Template(user, user, make([]*a.User, 0)))
	}
	var openChatId uint
	var chatOwnerId uint
	if openChatTemplate != nil {
		openChatId = openChatTemplate.ChatId
		chatOwnerId = openChatTemplate.Owner.UserId
	}
	userTemplate := user.Template(openChatId, chatOwnerId, user.Id)
	home := template.HomeTemplate{
		Chats:         chatTemplates,
		OpenChat:      openChatTemplate,
		User:          *userTemplate,
		IsAuthorized:  true,
		LoginTemplate: template.AuthTemplate{},
		Avatar:        avatarTmpl,
	}
	return home.HTML()
}
