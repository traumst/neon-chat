package handler

import (
	"log"
	"neon-chat/src/db"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/state"
	ti "neon-chat/src/interface"
	"neon-chat/src/model/app"
)

func TemplateOpenChat(state *state.State, db *db.DBConn, user *app.User) ti.Renderable {
	openChatId := state.GetOpenChat(user.Id)
	if openChatId == 0 {
		log.Printf("templateOpenchat DEBUG, user[%d] has no open chat\n", user.Id)
		return nil
	}
	openChat, err := shared.GetChat(state, db, user, openChatId)
	if err != nil {
		log.Printf("templateOpenchat ERROR, failed to get chat[%d], %s\n", openChatId, err.Error())
		return nil // TODO custom error pop-up
	}
	appChatUsers, err := shared.GetChatUsers(db, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] users, %s\n", openChatId, err.Error())
		return nil
	}
	appMsgs, err := shared.GetChatMessages(db, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] messages, %s\n", openChatId, err.Error())
		return nil
	}
	return openChat.Template(user, user, appChatUsers, appMsgs)
}
