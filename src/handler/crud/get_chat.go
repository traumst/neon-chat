package crud

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
	"neon-chat/src/utils"
)

func GetChat(state *state.State, dbConn sqlx.Ext, user *app.User, chatId uint) (*app.Chat, error) {
	dbChat, err := db.GetChat(dbConn, chatId)
	if err != nil {
		return nil, fmt.Errorf("chat[%d] not found in db: %s", chatId, err.Error())
	}
	if dbChat == nil {
		return nil, nil
	}
	chatIds, err := db.GetUserChatIds(dbConn, user.Id)
	if err != nil {
		return nil, fmt.Errorf("error getting chat ids for user[%d]: %s", user.Id, err.Error())
	}
	if !utils.Contains(chatIds, chatId) {
		return nil, fmt.Errorf("user[%d] is not in chat[%d]", user.Id, chatId)
	}
	dbOwner, err := db.GetUser(dbConn, dbChat.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat owner[%d]: %s", dbChat.OwnerId, err.Error())
	}
	return convert.ChatDBToApp(dbChat, dbOwner), nil
}
