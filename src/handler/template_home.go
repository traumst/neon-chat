package handler

import (
	"log"
	"net/http"
	"prplchat/src/convert"
	"prplchat/src/db"
	"prplchat/src/handler/state"
	ti "prplchat/src/interface"
	"prplchat/src/model/app"
	t "prplchat/src/model/template"
	h "prplchat/src/utils/http"
)

func TemplateHome(
	state *state.State,
	db *db.DBConn,
	r *http.Request,
	user *app.User,
) (string, error) {
	var avatarTmpl t.AvatarTemplate
	if dbAvatar, err := db.GetAvatar(user.Id); dbAvatar != nil && err == nil {
		avatar := convert.AvatarDBToApp(dbAvatar)
		avatarTmpl = avatar.Template(user)
	}
	openChatTemplate := TemplateOpenChat(state, db, user)
	chats, err := GetChats(state, db, user.Id)
	if err != nil {
		log.Printf("[%s] templateHome ERROR, failed getting chats for user[%d], %s\n",
			h.GetReqId(r), user.Id, err.Error())
		return "", err
	}
	var chatTemplates []ti.Renderable
	for _, chat := range chats {
		chatTemplates = append(chatTemplates, chat.Template(user, user, nil, nil))
	}
	var openChatId uint
	var chatOwnerId uint
	if openChatTemplate != nil {
		openChatId = openChatTemplate.(t.ChatTemplate).ChatId
		chatOwnerId = openChatTemplate.(t.ChatTemplate).Owner.(t.UserTemplate).UserId
	}
	userTemplate := user.Template(openChatId, chatOwnerId, user.Id)
	home := t.HomeTemplate{
		Chats:         chatTemplates,
		OpenChat:      openChatTemplate,
		User:          userTemplate,
		IsAuthorized:  true,
		LoginTemplate: t.AuthTemplate{},
		Avatar:        avatarTmpl,
	}
	return home.HTML()
}
