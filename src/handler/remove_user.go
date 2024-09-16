package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func RemoveUser(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
	expelledId uint,
) (*app.User, error) {
	log.Printf("ExpelUser TRACE user[%d] expells[%d] from chat[%d]\n", user.Id, expelledId, chatId)
	bothCanChat, err := db.UsersCanChat(dbConn.Conn, chatId, user.Id, expelledId)
	if err != nil {
		return nil, fmt.Errorf("failed to verify users can chat, %s", err.Error())
	} else if !bothCanChat {
		return nil, fmt.Errorf("at least one of users can't chat, activeUser[%d], expelled[%d]", user.Id, expelledId)
	}
	// veryfy user can only either leave themselves or be expelled by the owner
	if user.Id != expelledId {
		chat, err := GetChat(state, dbConn.Conn, user, chatId)
		if err != nil {
			log.Printf("ExpelUser ERROR user[%d] cannot find chat[%d], %s\n", user.Id, chatId, err.Error())
			return nil, fmt.Errorf("user cannot find chat, %s", err.Error())
		}
		if user.Id != chat.OwnerId {
			log.Printf("ExpelUser ERROR user[%d] cannot expel user[%d] from chat[%d]\n", user.Id, expelledId, chatId)
			return nil, fmt.Errorf("failed to expel user from chat")
		}
	}
	dbExpelled, err := db.GetUser(dbConn.Conn, expelledId)
	if err != nil || dbExpelled == nil {
		return nil, fmt.Errorf("user[%d] not found in db", expelledId)
	}
	log.Printf("ExpelUser TRACE removing[%d] from chat[%d]\n", expelledId, chatId)

	if dbConn.Tx == nil {
		log.Printf("ExpelUser ERROR no transaction provided\n")
		return nil, fmt.Errorf("no transaction provided")
	}

	err = db.RemoveChatUser(dbConn.Tx, chatId, expelledId)
	if err != nil {
		return nil, fmt.Errorf("failed to remove user[%d] from chat[%d]: %s", expelledId, chatId, err.Error())
	}
	log.Printf("ExpelUser TRACE closing chat[%d]\n", chatId)
	err = state.CloseChat(expelledId, chatId)
	if err != nil {
		log.Printf("ExpelUser TRACE user[%d] did not have chat[%d] open: %s", expelledId, chatId, err.Error())
		return nil, fmt.Errorf("failed to close chat[%d]", chatId)
	}
	log.Printf("ExpelUser TRACE user[%d] has been expelled from chat[%d]\n", expelledId, chatId)
	appExpelled := convert.UserDBToApp(dbExpelled, nil)
	return appExpelled, nil
}
