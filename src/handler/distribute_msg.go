package handler

import (
	"fmt"
	"log"
	"sync"

	"go.chat/src/model/app"
	"go.chat/src/model/event"
)

func DistributeMsg(
	state *AppState,
	chat *app.Chat,
	authorId uint,
	msg *app.Message,
	evnt event.UpdateType,
) error {
	users, err := chat.GetUsers(authorId)
	if err != nil || users == nil {
		return fmt.Errorf("DistributeMsg: get users, chat[%+v], %s", chat, err)
	}
	if len(users) == 0 {
		return fmt.Errorf("DistributeMsg: chatUsers are empty, chat[%+v], %s", chat, err)
	}

	var wg sync.WaitGroup
	var errors []string
	for _, user := range users {
		if user.Id == authorId {
			log.Printf("∞----> DistributeMsg TRACE new message is not sent to author[%d]\n", user.Id)
			continue
		}
		wg.Add(1)
		go func(user app.User, msg app.Message) {
			defer wg.Done()
			log.Printf("∞----> DistributeMsg TRACE new message will be sent to user[%d]\n", user.Id)
			data, err := msg.Template(&user).HTML()
			if err != nil {
				errors = append(errors, err.Error())
				return
			}
			err = distributeMsgToUser(state, chat.Id, msg.Id, user.Id, authorId, evnt, data)
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
	evnt event.UpdateType,
	data string,
) error {
	log.Printf("∞----> distributeMsgToUser TRACE user[%d] chat[%d] event[%v]\n", userId, chatId, evnt)
	openChat := state.GetOpenChat(userId)
	if openChat == nil {
		log.Printf("<----- distributeMsgToUser INFO user[%d] has no open chat to distribute", userId)
		return nil
	}
	if openChat.Id != chatId {
		log.Printf("<----- distributeMsgToUser INFO user[%d] has open chat[%d] different from message chat[%d]",
			userId, openChat.Id, chatId)
		return nil
	}

	conn, err := state.GetConn(userId)
	if err != nil {
		log.Printf("<----- distributeMsgToUser INFO user[%d] not connected, err:%s", userId, err.Error())
		return nil
	}

	msg := event.LiveUpdate{
		Event:    evnt,
		ChatId:   chatId,
		MsgId:    msgId,
		AuthorId: authorId,
		UserId:   userId,
	}

	switch evnt {
	case event.MessageAdded:
		msg.Data = data
		conn.In <- msg
		return nil
	case event.MessageDeleted:
		msg.Data = "[deletedM]"
		conn.In <- msg
		return nil
	default:
		return fmt.Errorf("unknown event type: %v", evnt)
	}
}
