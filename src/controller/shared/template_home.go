package shared

import (
	"fmt"
	"neon-chat/src/app"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/handler/pub"
	ti "neon-chat/src/interfaces"
	"neon-chat/src/state"
	"neon-chat/src/template"
)

func TemplateHome(state *state.State, dbConn *db.DBConn, user *app.User) (string, error) {
	var avatarTmpl template.AvatarTemplate
	if dbAvatar, err := db.GetAvatar(dbConn.Conn, user.Id); dbAvatar != nil && err == nil {
		avatar := convert.AvatarDBToApp(dbAvatar)
		avatarTmpl = avatar.Template(user)
	}
	openChatTemplate := TemplateOpenChat(state, dbConn, user)
	chats, err := pub.GetChats(dbConn.Conn, user.Id)
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
