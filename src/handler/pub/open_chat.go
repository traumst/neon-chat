package pub

import (
	"fmt"
	"log"

	"neon-chat/src/app"
	"neon-chat/src/db"
	"neon-chat/src/handler/priv"
	"neon-chat/src/state"
)

func OpenChat(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
) (*app.Chat, error) {
	err := state.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		return nil, fmt.Errorf("opening chat[%d] for user[%d], %s", chatId, user.Id, err.Error())
	}
	openChatId := state.GetOpenChat(user.Id)
	if openChatId != chatId {
		panic(fmt.Errorf("chat[%d] should have been open for user[%d]", chatId, user.Id))
	}
	appChat, err := priv.GetChat(state, dbConn.Conn, user, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		return nil, fmt.Errorf("getting chat[%d] for user[%d], %s", chatId, user.Id, err.Error())
	}
	return appChat, nil
}
