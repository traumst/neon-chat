package handler

import (
	"fmt"
	"log"
	d "neon-chat/src/db"
	"neon-chat/src/model/app"
	a "neon-chat/src/model/app"
	"neon-chat/src/state"
)

func LeaveChat(state *state.State, db *d.DBConn, user *a.User, chatId uint) (*app.Chat, *app.User, error) {
	chat, err := GetChat(state, db.Conn, user, chatId)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to leave chat: %s", err.Error())
	}
	if user.Id == chat.OwnerId {
		log.Printf("HandleUserLeaveChat ERROR cannot leave chat[%d] as owner\n", chatId)
		return nil, nil, fmt.Errorf("creator cannot leave chat")
	}
	log.Printf("HandleUserLeaveChat TRACE user[%d] leaves chat[%d]\n", user.Id, chatId)
	expelled, err := removeUser(state, db, user, chatId, user.Id)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR user[%d] failed to leave chat[%d], %s\n", user.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to leave from chat: %s", err.Error())
	}
	return chat, expelled, nil
}
