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
	targetUser string, // who to inform, empty for all users in chat
	subjectUser string, // which user changed
	event e.UpdateType,
) error {
	var targetUsers []string
	var err error
	if targetUser != "" {
		targetUsers = []string{targetUser}
	} else {
		targetUsers, err = chat.GetUsers(author)
		if err != nil {
			err = fmt.Errorf("DistributeChat: get users, chat[%d], %s", chat.ID, err)
		} else if len(targetUsers) == 0 {
			err = fmt.Errorf("DistributeChat: chatUsers are empty, chat[%+v], %s", chat, err)
		}
	}

	if err != nil {
		log.Printf("∞----> DistributeChat ERROR, %s\n", err)
		return err
	}
	targetUsersCount := len(targetUsers)
	if targetUsersCount < 1 {
		log.Printf("∞----> DistributeChat WARN targetUsersCount[%d] < 1\n", targetUsersCount)
		return nil
	}

	var wg sync.WaitGroup
	var errors []string
	wg.Add(targetUsersCount)
	for _, user := range targetUsers {
		go func() {
			defer wg.Done()
			log.Printf("∞----> DistributeChat TRACE event[%v] about subject[%s] will be sent to user[%s] in chat[%d]\n",
				event, subjectUser, user, chat.ID)
			err := distributeChatToUser(
				state,
				author,
				user,
				chat,
				subjectUser,
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
		log.Printf("∞----> DistributeChat ERROR occurred during distribution: %v\n", errors)
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
	subjectUser string,
	event e.UpdateType,
) error {
	conn, err := state.GetConn(targetUser)
	if err != nil {
		return err
	}
	if conn.User != targetUser {
		return fmt.Errorf("user[%s] does not own conn[%v], user[%s] does", targetUser, conn.Origin, conn.User)
	}

	switch event {
	case e.ChatCreated:
		return chatCreate(conn, event, author, targetChat, subjectUser)
	case e.ChatInvite:
		return chatInvite(conn, event, author, targetChat, subjectUser)
	case e.ChatExpel:
		return chatExpel(conn, event, author, targetChat.ID, subjectUser)
	case e.ChatDeleted:
		return chatDelete(conn, event, author, targetChat.ID, subjectUser)
	case e.ChatClose:
		return chatClose(conn, event, author, targetChat.ID, subjectUser)
	default:
		return fmt.Errorf("unknown event type[%v]", event)
	}
}

func chatCreate(conn *model.Conn, event e.UpdateType, author string, targetChat *app.Chat, subject string) error {
	log.Printf("∞----> distributeChatToUser TRACE author[%s] created chat[%d], target[%s], subject[%s]\n",
		author, targetChat.ID, conn.User, subject)
	if author != conn.User || conn.User != subject {
		return fmt.Errorf("chat_create expects author, target and subject to be the same")
	}
	template := targetChat.Template(author)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- e.LiveUpdate{
		Event:  event,
		ChatID: targetChat.ID,
		UserID: author,
		MsgID:  -1,
		Author: author,
		Data:   data,
	}
	return nil
}

func chatInvite(conn *model.Conn, event e.UpdateType, author string, targetChat *app.Chat, subject string) error {
	log.Printf("∞----> distributeChatToUser TRACE author[%s] invited subject[%s] to chat[%d], target[%s]\n",
		author, subject, targetChat.ID, conn.User)
	if author == conn.User || conn.User != subject {
		return fmt.Errorf("chat_invite expects author to be diff from target, target same as subject")
	}
	template := targetChat.Template(subject)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- e.LiveUpdate{
		Event:  event,
		ChatID: targetChat.ID,
		UserID: subject,
		MsgID:  -1,
		Author: author,
		Data:   data,
	}
	return nil
}

func chatExpel(conn *model.Conn, event e.UpdateType, author string, chatID int, subject string) error {
	log.Printf("∞----> chatExpel TRACE to user[%s] about author[%s] dropped subject[%s] from chat[%d]\n",
		conn.User, author, subject, chatID)
	if conn.User == author {
		return nil
	}
	conn.In <- e.LiveUpdate{
		Event:  event,
		ChatID: chatID,
		UserID: subject,
		MsgID:  -1,
		Author: author,
		Data:   "[expelU]",
	}
	return nil
}

func chatDelete(conn *model.Conn, event e.UpdateType, author string, chatID int, subject string) error {
	log.Printf("∞----> chatDelete TRACE author[%s] deleted chat[%d] for subject[%s], target[%s]\n",
		author, chatID, subject, conn.User)
	if subject != "" && conn.User != subject {
		return fmt.Errorf("chat_delete expect target and subject to be the same")
	}
	conn.In <- e.LiveUpdate{
		Event:  event,
		ChatID: chatID,
		UserID: conn.User,
		MsgID:  -1,
		Author: author,
		Data:   "[deletedC]",
	}
	return nil
}

func chatClose(conn *model.Conn, event e.UpdateType, author string, chatID int, subject string) error {
	log.Printf("∞----> chatClose TRACE user[%s] closed chat[%d] for target[%s], subject[%s]\n",
		author, chatID, conn.User, subject)
	if subject != "" && conn.User != subject {
		return fmt.Errorf("chat_close expect target and subject to be the same")
	}
	welcome := template.WelcomeTemplate{ActiveUser: conn.User}
	data, err := welcome.HTML()
	if err != nil {
		return err
	}
	conn.In <- e.LiveUpdate{
		Event:  event,
		ChatID: chatID,
		UserID: conn.User,
		MsgID:  -1,
		Author: author,
		Data:   data,
	}
	return nil
}
