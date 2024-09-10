package handler

import (
	"fmt"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/state"
	ti "neon-chat/src/interface"
	a "neon-chat/src/model/app"
	t "neon-chat/src/model/template"
)

func TemplateHome(state *state.State, db *d.DBConn, user *a.User) (string, error) {
	var avatarTmpl t.AvatarTemplate
	if dbAvatar, err := db.GetAvatar(user.Id); dbAvatar != nil && err == nil {
		avatar := convert.AvatarDBToApp(dbAvatar)
		avatarTmpl = avatar.Template(user)
	}
	openChatTemplate := TemplateOpenChat(state, db, user)
	chats, err := shared.GetChats(db, user.Id)
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
