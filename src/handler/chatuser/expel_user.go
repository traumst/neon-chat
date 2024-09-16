package chatuser

import (
	"fmt"
	"log"
	"neon-chat/src/db"
	"neon-chat/src/handler/chat"
	"neon-chat/src/model/app"
	"neon-chat/src/model/event"
	"neon-chat/src/sse"
	"neon-chat/src/state"
)

func ExpelUser(state *state.State, dbConn *db.DBConn, user *app.User, chatId uint, expelledId uint) (*app.User, error) {
	appExpelled, err := removeUser(state, dbConn, user, chatId, uint(expelledId))
	if err != nil {
		log.Printf("HandleUserExpelled ERROR failed to expell, %s\n", err.Error())
		return nil, fmt.Errorf("failed to expell user, %s", err.Error())
	}
	targetChat, err := chat.GetChat(state, dbConn.Tx, user, chatId)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return nil, fmt.Errorf("failed to expell user: %s", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, targetChat, user, appExpelled, appExpelled, event.ChatClose)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat close, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat close: %s", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, targetChat, user, appExpelled, appExpelled, event.ChatDrop)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat deleted, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat deleted: %s", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, targetChat, user, nil, appExpelled, event.ChatExpel)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat expel, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat expel: %s", err.Error())
	}
	return appExpelled, nil
}
