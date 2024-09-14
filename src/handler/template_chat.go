package handler

import (
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/state"
	ti "neon-chat/src/interface"
	"neon-chat/src/model/app"
)

func TemplateOpenChat(state *state.State, dbConn *db.DBConn, user *app.User) ti.Renderable {
	openChatId := state.GetOpenChat(user.Id)
	if openChatId == 0 {
		log.Printf("templateOpenchat DEBUG, user[%d] has no open chat\n", user.Id)
		return nil
	}
	openChat, err := shared.GetChat(state, dbConn.Conn, user, openChatId)
	if err != nil {
		log.Printf("templateOpenchat ERROR, failed to get chat[%d], %s\n", openChatId, err.Error())
		return nil // TODO custom error pop-up
	}
	appChatUsers, err := shared.GetChatUsers(dbConn.Conn, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] users, %s\n", openChatId, err.Error())
		return nil
	}
	appMsgs, err := shared.GetChatMessages(dbConn.Conn, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] messages, %s\n", openChatId, err.Error())
		return nil
	}
	if user.Avatar == nil {
		dbAvatar, err := db.GetAvatar(dbConn.Conn, user.Id)
		if dbAvatar != nil && err == nil {
			user.Avatar = convert.AvatarDBToApp(dbAvatar)
		}
	}
	return openChat.Template(user, user, appChatUsers, appMsgs)
}
