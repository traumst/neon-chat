package pub

import (
	"fmt"
	"log"

	"neon-chat/src/app"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/state"
)

func DeleteChat(state *state.State, dbConn *db.DBConn, user *app.User, chatId uint) (*app.Chat, error) {
	dbChat, err := db.GetChat(dbConn.Tx, chatId)
	if err != nil {
		log.Printf("INFO DeleteChat chat[%d] not found in db: %s", chatId, err)
		return nil, nil
	}
	if user.Id != dbChat.OwnerId {
		return nil, fmt.Errorf("user[%d] is not owner[%d] of chat[%d]", user.Id, dbChat.OwnerId, chatId)
	}
	err = db.DeleteChat(dbConn.Tx, chatId)
	if err != nil {
		return nil, fmt.Errorf("error deleting chat in db: %s", err.Error())
	}
	chat := convert.ChatDBToApp(dbChat, &db.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	})
	return chat, nil
}
