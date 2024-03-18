package handler

import (
	"fmt"
	"log"
	"sync"

	"go.chat/model"
	"go.chat/model/app"
	e "go.chat/model/event"
	"go.chat/model/template"
)

// empty targetUser means all users in chat
func DistributeChat(
	state *model.AppState,
	chat *app.Chat,
	author string, // who made the change
	targetUser string, // who to inform
	event e.UpdateType,
) error {
	var err error

	targetUsers := []string{targetUser}
	if targetUser == "" {
		targetUsers, err = chat.GetUsers(author)
		if err != nil || targetUsers == nil {
			return fmt.Errorf("DistributeChat: get users, chat[%+v], %s", chat, err)
		}
		if len(targetUsers) == 0 {
			return fmt.Errorf("DistributeChat: chatUsers are empty, chat[%+v], %s", chat, err)
		}
	}

	var wg sync.WaitGroup
	var errors []string
	for _, user := range targetUsers {
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
	targetChat *app.Chat,
	event e.UpdateType,
) error {
	conn, err := state.GetConn(targetUser)
	if err != nil {
		return err
	}
	if conn.User != targetUser {
		return fmt.Errorf("user[%s] does not own conn[%v], user[%s] does", targetUser, conn.Origin, conn.User)
	}

	var data string
	switch event {
	case e.ChatCreated:
		template := targetChat.Template(targetUser)
		data, err = template.ShortHTML()
		if err != nil {
			return err
		}
		conn.In <- e.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			UserID: targetUser,
			MsgID:  -1,
			Author: author,
			Data:   data,
		}
	case e.ChatInvite:
		template := targetChat.Template(targetUser)
		data, err = template.ShortHTML()
		if err != nil {
			return err
		}
		conn.In <- e.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			UserID: targetUser,
			MsgID:  -1,
			Author: author,
			Data:   data,
		}
	case e.ChatDeleted:
		log.Printf("∞----> distributeChatToUser TRACE user[%s] deleted chat[%d]\n",
			author, targetChat.ID)
		conn.In <- e.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			UserID: targetUser,
			MsgID:  -1,
			Author: author,
			Data:   "[deletedC]",
		}
	case e.ChatClose:
		welcome := template.WelcomeTemplate{ActiveUser: targetUser}
		data, err = welcome.HTML()
		if err != nil {
			return err
		}
		conn.In <- e.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			UserID: targetUser,
			MsgID:  -1,
			Author: author,
			Data:   data,
		}
	case e.ChatUserDrop:
		conn.In <- e.LiveUpdate{
			Event:  event,
			ChatID: targetChat.ID,
			UserID: targetUser,
			MsgID:  -1,
			Author: author,
			Data:   "[deletedU]",
		}
	default:
		return fmt.Errorf("unknown event type[%v]", event)
	}
	return nil
}
