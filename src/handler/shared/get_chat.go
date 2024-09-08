package shared

import (
	"fmt"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler/state"
	a "neon-chat/src/model/app"
	"neon-chat/src/utils"
)

func GetChat(state *state.State, db *d.DBConn, user *a.User, chatId uint) (*a.Chat, error) {
	dbChat, err := db.GetChat(chatId)
	if err != nil {
		return nil, fmt.Errorf("chat[%d] not found in db: %s", chatId, err.Error())
	}
	if dbChat == nil {
		return nil, nil
	}
	chatIds, err := db.GetUserChatIds(user.Id)
	if err != nil {
		return nil, fmt.Errorf("error getting chat ids for user[%d]: %s", user.Id, err.Error())
	}
	if !utils.Contains(chatIds, chatId) {
		return nil, fmt.Errorf("user[%d] is not in chat[%d]", user.Id, chatId)
	}
	dbOwner, err := db.GetUser(dbChat.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat owner[%d]: %s", dbChat.OwnerId, err.Error())
	}
	return convert.ChatDBToApp(dbChat, dbOwner), nil
}
