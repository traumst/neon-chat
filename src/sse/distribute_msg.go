package sse

import (
	"fmt"
	"log"
	"sync"

	"neon-chat/src/app"
	"neon-chat/src/event"
	"neon-chat/src/handler/pub"
	"neon-chat/src/state"

	"github.com/jmoiron/sqlx"
)

func DistributeMsg(
	state *state.State,
	dbConn sqlx.Ext,
	chat *app.Chat,
	msg *app.Message,
	updateType event.EventType,
) error {
	if chat == nil || msg == nil {
		return fmt.Errorf("mandatory argument/s cannot be nil")
	}
	users, err := pub.GetChatUsers(dbConn, chat.Id)
	if err != nil {
		return fmt.Errorf("failed to get users in chat[%d], %s", chat.Id, err)
	}
	if len(users) <= 0 {
		return fmt.Errorf("chatUsers are empty, chat[%d], %s", chat.Id, err)
	}
	var owner *app.User
	for _, user := range users {
		if owner == nil && user.Id == chat.OwnerId {
			owner = user
			break
		}
	}

	var errors []string
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(viewer app.User, msg app.Message) {
			defer wg.Done()
			log.Printf("DistributeMsg TRACE event[%s] will be sent to user[%d]\n", updateType, viewer.Id)
			msgTmpl, err := msg.Template(&viewer, owner, msg.Quote)
			if err != nil {
				errors = append(errors, fmt.Sprintf("template:%s", err.Error()))
				return
			}
			msgData, err := msgTmpl.HTML()
			if err != nil {
				errors = append(errors, fmt.Sprintf("html:%s", err.Error()))
				return
			}
			err = distributeMsgToUser(state, chat.Id, msg.Id, viewer.Id, msg.Author.Id, updateType, msgData)
			if err != nil {
				errors = append(errors, fmt.Sprintf("distribute:%s", err.Error()))
			}
		}(*user, *msg)
	}

	wg.Wait()
	if len(errors) > 0 {
		return error(fmt.Errorf("DistributeMsg errors: [%v]", errors))
	} else {
		return nil
	}
}

func distributeMsgToUser(
	state *state.State,
	chatId uint,
	msgId uint,
	userId uint,
	authorId uint,
	updateType event.EventType,
	data string,
) error {
	log.Printf("distributeMsgToUser TRACE user[%d] chat[%d] event[%v]\n", userId, chatId, updateType)
	openChatId := state.GetOpenChat(userId)
	if openChatId == 0 {
		log.Printf("distributeMsgToUser INFO user[%d] has no open chat to distribute", userId)
		return nil
	}
	// TODO only sends updates to open chat
	if openChatId != chatId {
		log.Printf("distributeMsgToUser INFO user[%d] has open chat[%d] different from message chat[%d]",
			userId, openChatId, chatId)
		return nil
	}
	msg := event.LiveEvent{
		Event:    updateType,
		ChatId:   chatId,
		MsgId:    msgId,
		AuthorId: authorId,
		UserId:   userId,
	}
	switch updateType {
	case event.MessageAdd:
		msg.Data = data
	case event.MessageDrop:
		msg.Data = "[deletedM]"
	default:
		return fmt.Errorf("unknown event type: %v", updateType)
	}
	conns := state.GetConn(userId)
	for _, conn := range conns {
		conn.In <- msg
	}
	return nil
}
