package handler

import (
	"fmt"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	ti "neon-chat/src/interface"
	"neon-chat/src/model/app"
	"neon-chat/src/model/template"
	"neon-chat/src/state"
)

func TemplateHome(state *state.State, dbConn *db.DBConn, user *app.User) (string, error) {
	var avatarTmpl template.AvatarTemplate
	if dbAvatar, err := db.GetAvatar(dbConn.Conn, user.Id); dbAvatar != nil && err == nil {
		avatar := convert.AvatarDBToApp(dbAvatar)
		avatarTmpl = avatar.Template(user)
	}
	openChatTemplate := TemplateOpenChat(state, dbConn, user)
	chats, err := GetChats(dbConn.Conn, user.Id)
	if err != nil {
		return "", fmt.Errorf("failed getting chats for user, %s", err)
	}
	var chatTemplates []ti.Renderable
	for _, chat := range chats {
		chatTemplates = append(chatTemplates, chat.Template(user, user, nil, nil))
	}
	var openChatId uint
	var chatOwnerId uint
	if openChatTemplate != nil {
		openChatId = openChatTemplate.(template.ChatTemplate).ChatId
		chatOwnerId = openChatTemplate.(template.ChatTemplate).Owner.(template.UserTemplate).UserId
	}
	userTemplate := user.Template(openChatId, chatOwnerId, user.Id)
	home := template.HomeTemplate{
		Chats:         chatTemplates,
		OpenChat:      openChatTemplate,
		User:          userTemplate,
		IsAuthorized:  true,
		LoginTemplate: template.AuthTemplate{},
		Avatar:        avatarTmpl,
	}
	return home.HTML()
}
