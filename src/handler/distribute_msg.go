package handler

import (
	"fmt"
	"log"
	"sync"

	"prplchat/src/model/app"
	"prplchat/src/model/event"
)

func DistributeMsg(
	state *AppState,
	chat *app.Chat,
	authorId uint,
	msg *app.Message,
	updateType event.UpdateType,
) error {
	// have to get users by owner - author may have been removed
	users, err := chat.GetUsers(chat.Owner.Id)
	if err != nil || users == nil {
		return fmt.Errorf("DistributeMsg: get users, chat[%d], %s", chat.Id, err)
	}
	if len(users) <= 0 {
		return fmt.Errorf("DistributeMsg: chatUsers are empty, chat[%d], %s", chat.Id, err)
	}

	var wg sync.WaitGroup
	var errors []string
	for _, user := range users {
		if user.Id == authorId {
			log.Printf("DistributeMsg TRACE new message is not sent to author[%d]\n", user.Id)
			continue
		}
		wg.Add(1)
		go func(user app.User, msg app.Message) {
			defer wg.Done()
			log.Printf("DistributeMsg TRACE new message will be sent to user[%d]\n", user.Id)
			data, err := msg.Template(&user).HTML()
			if err != nil {
				errors = append(errors, err.Error())
				return
			}
			err = distributeMsgToUser(state, chat.Id, msg.Id, user.Id, authorId, updateType, data)
			if err != nil {
				errors = append(errors, err.Error())
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
	state *AppState,
	chatId int,
	msgId int,
	userId uint,
	authorId uint,
	updateType event.UpdateType,
	data string,
) error {
	log.Printf("distributeMsgToUser TRACE user[%d] chat[%d] event[%v]\n", userId, chatId, updateType)
	openChat := state.GetOpenChat(userId)
	if openChat == nil {
		log.Printf("distributeMsgToUser INFO user[%d] has no open chat to distribute", userId)
		return nil
	}
	if openChat.Id != chatId {
		log.Printf("distributeMsgToUser INFO user[%d] has open chat[%d] different from message chat[%d]",
			userId, openChat.Id, chatId)
		return nil
	}

	conn, err := state.GetConn(userId)
	if err != nil {
		log.Printf("distributeMsgToUser INFO user[%d] not connected, err:%s", userId, err.Error())
		return nil
	}

	msg := event.LiveUpdate{
		Event:    updateType,
		ChatId:   chatId,
		MsgId:    msgId,
		AuthorId: authorId,
		UserId:   userId,
	}

	switch updateType {
	case event.MessageAdd:
		msg.Data = data
		conn.In <- msg
		return nil
	case event.MessageDrop:
		msg.Data = "[deletedM]"
		conn.In <- msg
		return nil
	default:
		return fmt.Errorf("unknown event type: %v", updateType)
	}
}
