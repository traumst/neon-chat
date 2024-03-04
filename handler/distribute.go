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
		return fmt.Errorf("get users, chat[%+v], %s", chat, err)
	}
	if len(users) == 0 {
		return fmt.Errorf("chatUsers are empty, chat[%+v], %s", chat, err)
	}

	var wg sync.WaitGroup
	var errors []string
	for _, user := range users {
		if user == author {
			log.Printf("∞----> distributeBetween TRACE new message is not sent to author[%s]\n", user)
			continue
		}
		wg.Add(1)
		go func(user string, msg model.Message) {
			defer wg.Done()
			data, err := msg.ToTemplate(user).GetHTML()
			if err != nil {
				errors = append(errors, err.Error())
				return
			}
			err = distributeMsgToUser(state, chat.ID, user, author, event, data)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}(user, *msg)
	}

	wg.Wait()
	if len(errors) > 0 {
		return error(fmt.Errorf("distributing message errors [%v]", errors))
	} else {
		return nil
	}
}

func DistributeChat(
	state *model.AppState,
	user string,
	template *model.ChatTemplate,
	event model.UpdateType,
) error {
	shortHtml, err := template.GetShortHTML()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = distributeChatToUser(
			state,
			user,
			template.ID,
			template.User,
			shortHtml,
			event,
		)
	}()

	wg.Wait()
	return err
}

func distributeMsgToUser(
	state *model.AppState,
	chatID int,
	user string,
	author string,
	event model.UpdateType,
	data string,
) error {
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
			Author: author,
		}
		return nil
	default:
		return fmt.Errorf("unknown event type: %s", event.String())
	}
}

func distributeChatToUser(
	state *model.AppState,
	targetUser string,
	chatID int,
	chatAuthor string,
	html string,
	event model.UpdateType,
) error {
	conn, err := state.GetConn(targetUser)
	if err != nil {
		return fmt.Errorf("user[%s] not connected, err:%s", targetUser, err.Error())
	}
	if conn.User != targetUser {
		return fmt.Errorf("user[%s] does not own conn[%v], user[%s] does", targetUser, conn.Origin, conn.User)
	}

	switch event {
	case model.ChatCreated, model.ChatInvite, model.ChatDeleted:
		conn.In <- model.LiveUpdate{
			Event:  event,
			ChatID: chatID,
			Author: chatAuthor,
			Data:   html,
		}
		return nil
	default:
		return fmt.Errorf("unknown event type[%s]", event.String())
	}
}
