package crud

import (
	"fmt"
	"log"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func ExpelUser(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
	expelledId uint,
) (*app.Chat, *app.User, error) {
	appExpelled, err := removeUser(state, dbConn, user, chatId, uint(expelledId))
	if err != nil {
		log.Printf("HandleUserExpelled ERROR failed to expell, %s\n", err.Error())
		return nil, nil, fmt.Errorf("failed to expell user, %s", err.Error())
	}
	targetChat, err := GetChat(state, dbConn.Tx, user, chatId)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to expell user: %s", err.Error())
	}
	return targetChat, appExpelled, nil
}
