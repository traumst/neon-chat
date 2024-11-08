package pub

import (
	"fmt"
	"log"
	"neon-chat/src/app"
	"neon-chat/src/db"
	"neon-chat/src/handler/priv"
	"neon-chat/src/state"
)

func LeaveChat(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
) (*app.Chat, *app.User, error) {
	chat, err := priv.GetChat(state, dbConn.Conn, user, chatId)
	if err != nil {
		log.Printf("ERROR HandleUserLeaveChat cannot find chat[%d], %s\n", chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to leave chat: %s", err.Error())
	}
	if user.Id == chat.OwnerId {
		log.Printf("ERROR HandleUserLeaveChat cannot leave chat[%d] as owner\n", chatId)
		return nil, nil, fmt.Errorf("creator cannot leave chat")
	}
	log.Printf("TRACE HandleUserLeaveChat user[%d] leaves chat[%d]\n", user.Id, chatId)
	expelled, err := priv.RemoveUser(state, dbConn, user, chatId, user.Id)
	if err != nil {
		log.Printf("ERROR HandleUserLeaveChat user[%d] failed to leave chat[%d], %s\n", user.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to leave from chat: %s", err.Error())
	}
	return chat, expelled, nil
}
