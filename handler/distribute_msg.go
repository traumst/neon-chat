package handler

import (
	"fmt"
	"log"
	"sync"

	"go.chat/model"
)

func DistributeMsg(
	state *model.AppState,
	chat *model.Chat,
	author string,
	msg *model.Message,
	event model.UpdateType,
) error {
	users, err := chat.GetUsers(author)
	if err != nil || users == nil {
		return fmt.Errorf("DistributeMsg: get users, chat[%+v], %s", chat, err)
	}
	if len(users) == 0 {
		return fmt.Errorf("DistributeMsg: chatUsers are empty, chat[%+v], %s", chat, err)
	}

	var wg sync.WaitGroup
	var errors []string
	for _, user := range users {
		if user == author {
			log.Printf("∞----> DistributeMsg TRACE new message is not sent to author[%s]\n", user)
			continue
		}
		wg.Add(1)
		go func(user string, msg model.Message) {
			defer wg.Done()
			log.Printf("∞----> DistributeMsg TRACE new message will be sent to user[%s]\n", user)
			data, err := msg.ToTemplate(user).GetHTML()
			if err != nil {
				errors = append(errors, err.Error())
				return
			}
			err = distributeMsgToUser(state, chat.ID, msg.ID, user, author, event, data)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}(user, *msg)
	}

	wg.Wait()
	if len(errors) > 0 {
		return error(fmt.Errorf("DistributeMsg errors: [%v]", errors))
	} else {
		return nil
	}
}

func distributeMsgToUser(
	state *model.AppState,
	chatID int,
	msgID int,
	user string,
	author string,
	event model.UpdateType,
	data string,
) error {
	log.Printf("∞----> distributeMsgToUser TRACE user[%s] chat[%d] event[%v]\n", user, chatID, event)
	openChat := state.GetOpenChat(user)
	if openChat == nil {
		return fmt.Errorf("user[%s] has no open chat", user)
	}
	if openChat.ID != chatID {
		return fmt.Errorf("user[%s] has open chat[%d] different from message chat[%d]", user, openChat.ID, chatID)
	}

	conn, err := state.GetConn(user)
	if err != nil {
		return fmt.Errorf("user[%s] not connected, err:%s", user, err.Error())
	}

	switch event {
	case model.MessageAdded, model.MessageDeleted:
		conn.In <- model.LiveUpdate{
			Event:  event,
			Data:   data,
			ChatID: chatID,
			MsgID:  msgID,
			Author: author,
		}
		return nil
	default:
		return fmt.Errorf("unknown event type: %s", event.String())
	}
}
