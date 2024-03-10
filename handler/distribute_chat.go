package handler

import (
	"fmt"
	"log"
	"sync"

	"go.chat/model"
)

func DistributeChat(
	state *model.AppState,
	chat *model.Chat,
	author string, // who made the change
	targetUser string, // who to inform
	event model.UpdateType,
) error {
	var users []string
	var err error
	// AL TODO bad logic
	// 	if targetUser IS NOT author - take only targetUser
	// 	if targetUser IS author 	- take all users
	if targetUser != "" && targetUser != author {
		users = []string{targetUser}
	} else {
		users, err = chat.GetUsers(author)
		if err != nil || users == nil {
			return fmt.Errorf("DistributeChat: get users, chat[%+v], %s", chat, err)
		}
		if len(users) == 0 {
			return fmt.Errorf("DistributeChat: chatUsers are empty, chat[%+v], %s", chat, err)
		}
	}

	var wg sync.WaitGroup
	var errors []string
	for _, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("∞----> DistributeChat TRACE event[%v] will be sent to user[%s] in chat[%d]\n",
				event, user, chat.ID)
			err = distributeChatToUser(
				state,
				author,
				user,
				chat,
				event,
			)
			if err != nil {
				errors = append(errors, err.Error())
				return
			}
		}()
	}

	wg.Wait()
	if len(errors) > 0 {
		return error(fmt.Errorf("DistributeChat errors: [%v]", errors))
	} else {
		return nil
	}
}

func distributeChatToUser(
	state *model.AppState,
	author string,
	targetUser string,
	targetChat *model.Chat,
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
	case model.ChatCreated:
		template := targetChat.ToTemplate(targetUser)
		data, err = template.GetShortHTML()
		if err != nil {
			return err
		}
		conn.In <- model.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			MsgID:  -1,
			Author: author,
			Data:   data,
		}
	case model.ChatInvite:
		template := targetChat.ToTemplate(targetUser)
		data, err = template.GetShortHTML()
		if err != nil {
			return err
		}
		conn.In <- model.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			MsgID:  -1,
			Author: author,
			Data:   data,
		}
	case model.ChatDeleted:
		conn.In <- model.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			MsgID:  -1,
			Author: targetUser,
			Data:   "CHAT_DELETED_MESSAGE",
		}
	case model.ChatClose:
		welcome := model.WelcomeTemplate{ActiveUser: targetUser}
		data, err = welcome.GetHTML()
		if err != nil {
			return err
		}
		conn.In <- model.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			MsgID:  -1,
			Author: author,
			Data:   data,
		}
	default:
		return fmt.Errorf("unknown event type[%v]", event)
	}
	return nil
}
