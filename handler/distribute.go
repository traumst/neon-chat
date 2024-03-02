package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/model"
	"go.chat/utils"
)

func DistributeMsg(
	state *model.AppState,
	chat *model.Chat,
	author string,
	r *http.Request,
	event model.UpdateType,
	data string,
) error {
	users, err := chat.GetUsers(author)
	if err != nil || users == nil {
		return fmt.Errorf("get users, chat[%+v], %s", chat, err)
	}
	if len(users) == 0 {
		return fmt.Errorf("chatUsers are empty, chat[%+v], %s", chat, err)
	}

	reqId := utils.GetReqId(r)
	var wg sync.WaitGroup
	var errors []string
	for _, user := range users {
		if user == author {
			log.Printf("--%s-> distributeBetween TRACE new message is not sent to author[%s]\n", reqId, user)
			continue
		}
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			err := distributeMsgToUser(state, chat.ID, user, author, event, data)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}(user)
	}

	wg.Wait()
	if len(errors) > 0 {
		return error(fmt.Errorf("distributing message errors [%v]", errors))
	} else {
		return nil
	}
}

func DistributeChat(
	reqId string,
	state *model.AppState,
	user string,
	template *model.ChatTemplate,
	event model.UpdateType,
) error {
	log.Printf("--%s-> informOwner TRACE sending update of chat[%s] header to user [%s]\n", reqId, template.Name, user)
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

	var up model.LiveUpdate
	switch event {
	case model.MessageDeleted:
		up = model.LiveUpdate{
			Event:  event,
			ChatID: chatID,
		}
	case model.MessageAdded:
		up = model.LiveUpdate{
			Event:  event,
			Data:   data,
			ChatID: chatID,
			Author: author,
		}
	default:
		return fmt.Errorf("unknown event type: %s", event.String())
	}

	conn.In <- up
	return nil
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

	var up model.LiveUpdate
	switch event {
	case model.ChatDeleted:
		up = model.LiveUpdate{
			Event:  event,
			ChatID: chatID,
		}
	case model.ChatCreated, model.ChatInvite:
		up = model.LiveUpdate{
			Event:  event,
			ChatID: chatID,
			Author: chatAuthor,
			Data:   html,
		}
	default:
		return fmt.Errorf("unknown event type: %s", event.String())
	}

	conn.In <- up
	return nil
}
