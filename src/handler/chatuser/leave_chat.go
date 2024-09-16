package chatuser

import (
	"fmt"
	"log"
	d "neon-chat/src/db"
	c "neon-chat/src/handler/chat"
	a "neon-chat/src/model/app"
	"neon-chat/src/model/event"
	"neon-chat/src/sse"
	"neon-chat/src/state"
)

func LeaveChat(state *state.State, db *d.DBConn, user *a.User, chatId uint) error {
	chat, err := c.GetChat(state, db.Conn, user, chatId)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return fmt.Errorf("failed to leave chat: %s", err.Error())
	}
	if user.Id == chat.OwnerId {
		log.Printf("HandleUserLeaveChat ERROR cannot leave chat[%d] as owner\n", chatId)
		return fmt.Errorf("creator cannot leave chat")
	}
	log.Printf("HandleUserLeaveChat TRACE user[%d] leaves chat[%d]\n", user.Id, chatId)
	expelled, err := removeUser(state, db, user, chatId, user.Id)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR user[%d] failed to leave chat[%d], %s\n", user.Id, chatId, err.Error())
		return fmt.Errorf("failed to leave from chat: %s", err.Error())
	}
	log.Printf("HandleUserLeaveChat TRACE informing users in chat[%d]\n", chat.Id)
	err = sse.DistributeChat(state, db.Tx, chat, expelled, expelled, expelled, event.ChatClose)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat close, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, db.Tx, chat, expelled, expelled, expelled, event.ChatDrop)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat deleted, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, db.Tx, chat, expelled, nil, expelled, event.ChatLeave)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat user drop, %s\n", err.Error())
	}
	return nil
}
