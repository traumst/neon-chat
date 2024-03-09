package handler

import (
	"fmt"
	"sync"

	"go.chat/model"
)

func DistributeChat(
	state *model.AppState,
	user string,
	template *model.ChatTemplate,
	event model.UpdateType,
) error {
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = distributeChatToUser(
			state,
			user,
			template,
			event,
		)
	}()

	wg.Wait()
	return err
}

func distributeChatToUser(
	state *model.AppState,
	targetUser string,
	targetChat *model.ChatTemplate,
	event model.UpdateType,
) error {
	conn, err := state.GetConn(targetUser)
	if err != nil {
		return fmt.Errorf("user[%s] not connected, err:%s", targetUser, err.Error())
	}
	if conn.User != targetUser {
		return fmt.Errorf("user[%s] does not own conn[%v], user[%s] does", targetUser, conn.Origin, conn.User)
	}

	var data string
	switch event {
	case model.ChatCreated, model.ChatInvite:
		data, err = targetChat.GetShortHTML()
		if err != nil {
			return err
		}
		sendToChannel(&conn.In, event, targetChat.ChatID, targetChat.Owner, data)
		return nil
	case model.ChatDeleted:
		welcome := model.WelcomeTemplate{ActiveUser: targetUser}
		data, err = welcome.GetHTML()
		if err != nil {
			return err
		}
		sendToChannel(&conn.In, event, targetChat.ChatID, targetChat.Owner, data)
		return nil
	default:
		return fmt.Errorf("unknown event type[%s]", event.String())
	}
}

func sendToChannel(
	ch *chan model.LiveUpdate,
	event model.UpdateType,
	chatID int,
	targetUser string,
	data string,
) {
	*ch <- model.LiveUpdate{
		Event:  event,
		ChatID: chatID,
		MsgID:  -1,
		Author: targetUser,
		Data:   data,
	}
}
