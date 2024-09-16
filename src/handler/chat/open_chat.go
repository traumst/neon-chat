package chat

import (
	"fmt"
	"log"

	"neon-chat/src/db"
	m "neon-chat/src/handler/msg"
	u "neon-chat/src/handler/user"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func OpenChat(state *state.State, dbConn *db.DBConn, user *app.User, chatId uint) (string, error) {
	err := state.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		return TemplateWelcome(user)
	}
	openChatId := state.GetOpenChat(user.Id)
	if openChatId != chatId {
		panic(fmt.Errorf("chat[%d] should have been open for user[%d]", chatId, user.Id))
	}
	appChat, err := GetChat(state, dbConn.Conn, user, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		return TemplateWelcome(user)
	}
	appChatUsers, err := u.GetChatUsers(dbConn.Conn, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR getting chat[%d] users for user[%d], %s\n", chatId, user.Id, err.Error())
		return TemplateWelcome(user)
	}
	appChatMsgs, err := m.GetChatMessages(dbConn.Conn, chatId)
	if err != nil {
		log.Printf("HandleChatAdd ERROR getting chat[%d] messages: %s", appChat.Id, err)
	}
	return appChat.Template(user, user, appChatUsers, appChatMsgs).HTML()
}
