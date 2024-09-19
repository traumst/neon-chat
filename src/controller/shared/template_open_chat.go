package shared

import (
	"log"
	"neon-chat/src/db"
	"neon-chat/src/handler/priv"
	"neon-chat/src/handler/pub"
	ti "neon-chat/src/interfaces"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func TemplateOpenChat(state *state.State, dbConn *db.DBConn, user *app.User) ti.Renderable {
	openChatId := state.GetOpenChat(user.Id)
	if openChatId == 0 {
		log.Printf("templateOpenchat DEBUG, user[%d] has no open chat\n", user.Id)
		return nil
	}
	openChat, err := priv.GetChat(state, dbConn.Conn, user, openChatId)
	if err != nil {
		log.Printf("templateOpenchat ERROR, failed to get chat[%d], %s\n", openChatId, err.Error())
		return nil
	}
	appChatUsers, err := pub.GetChatUsers(dbConn.Conn, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] users, %s\n", openChatId, err.Error())
		return nil
	}
	appMsgs, err := pub.GetChatMessages(dbConn.Conn, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] messages, %s\n", openChatId, err.Error())
		return nil
	}
	if user.Avatar == nil {
		avatar, err := pub.GetAvatar(dbConn.Conn, user.Id)
		if avatar != nil && err == nil {
			user.Avatar = avatar
		}
	}
	return openChat.Template(user, user, appChatUsers, appMsgs)
}
