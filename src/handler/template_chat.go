package handler

import (
	"log"
	"prplchat/src/db"
	"prplchat/src/handler/state"
	"prplchat/src/model/app"
	"prplchat/src/model/template"
)

func TemplateOpenChat(state *state.State, db *db.DBConn, user *app.User) *template.ChatTemplate {
	openChatId := state.GetOpenChat(user.Id)
	if openChatId == 0 {
		log.Printf("templateOpenchat DEBUG, user[%d] has no open chat\n", user.Id)
		return nil
	}
	openChat, err := GetChat(state, db, user, openChatId)
	if err != nil {
		log.Printf("templateOpenchat ERROR, failed to get chat[%d], %s\n", openChatId, err.Error())
		return nil // TODO custom error pop-up
	}
	appChatUsers, err := GetChatUsers(db, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] users, %s\n", openChatId, err.Error())
		return nil
	}
	appMsgs, err := GetChatMessages(db, openChatId)
	if err != nil {
		log.Printf("templateOpenChat ERROR, failed getting chat[%d] messages, %s\n", openChatId, err.Error())
		return nil
	}
	return openChat.Template(user, user, appChatUsers, appMsgs)
}
