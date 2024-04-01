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
	author *app.User, // who made the change
	targetUser *app.User, // who to inform, empty for all users in chat
	subjectUser *app.User, // which user changed
	event e.UpdateType,
) error {
	var targetUsers []*app.User
	var err error
	if targetUser != nil {
		targetUsers = []*app.User{targetUser}
	} else {
		targetUsers, err = chat.GetUsers(author.Id)
		if err != nil {
			err = fmt.Errorf("DistributeChat: get users, chat[%d], %s", chat.Id, err)
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
			log.Printf("∞----> DistributeChat TRACE event[%v] about subject[%d] will be sent to user[%d] in chat[%d]\n",
				event, subjectUser.Id, user.Id, chat.Id)
			err := distributeChatToUser(
				state,
				author,
				user.Id,
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
	author *app.User,
	targetUserId uint,
	targetChat *app.Chat,
	subjectUser *app.User,
	event e.UpdateType,
) error {
	conn, err := state.GetConn(targetUserId)
	if err != nil {
		return err
	}
	if conn.User.Id != targetUserId {
		return fmt.Errorf("user[%d] does not own conn[%v], user[%d] does", targetUserId, conn.Origin, conn.User.Id)
	}

	switch event {
	case e.ChatCreated:
		return chatCreate(conn, event, author, targetChat, subjectUser.Id)
	case e.ChatInvite:
		return chatInvite(conn, event, author.Id, targetChat, subjectUser)
	case e.ChatExpel:
		return chatExpel(conn, event, author.Id, targetChat.Id, subjectUser.Id)
	case e.ChatDeleted:
		return chatDelete(conn, event, author.Id, targetChat.Id, subjectUser.Id)
	case e.ChatClose:
		return chatClose(conn, event, author.Id, targetChat.Id, subjectUser.Id)
	default:
		return fmt.Errorf("unknown event type[%v]", event)
	}
}

func chatCreate(conn *model.Conn, event e.UpdateType, author *app.User, targetChat *app.Chat, subjectId uint) error {
	log.Printf("∞----> distributeChatToUser TRACE author[%d] created chat[%d], target[%d], subject[%d]\n",
		author.Id, targetChat.Id, conn.User.Id, subjectId)
	if author.Id != conn.User.Id || conn.User.Id != subjectId {
		return fmt.Errorf("chat_create expects author, target and subject to be the same")
	}
	template := targetChat.Template(author)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- e.LiveUpdate{
		Event:    event,
		ChatId:   targetChat.Id,
		UserId:   author.Id,
		MsgId:    -1,
		AuthorId: author.Id,
		Data:     data,
	}
	return nil
}

func chatInvite(conn *model.Conn, event e.UpdateType, authorId uint, targetChat *app.Chat, subject *app.User) error {
	log.Printf("∞----> distributeChatToUser TRACE author[%d] invited subject[%d] to chat[%d], target[%d]\n",
		authorId, subject.Id, targetChat.Id, conn.User.Id)
	if authorId == conn.User.Id || conn.User.Id != subject.Id {
		return fmt.Errorf("chat_invite expects author to be diff from target, target same as subject")
	}
	template := targetChat.Template(subject)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- e.LiveUpdate{
		Event:    event,
		ChatId:   targetChat.Id,
		UserId:   subject.Id,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}

func chatExpel(conn *model.Conn, event e.UpdateType, authorId uint, chatId int, subjectId uint) error {
	log.Printf("∞----> chatExpel TRACE to user[%d] about author[%d] dropped subject[%d] from chat[%d]\n",
		conn.User.Id, authorId, subjectId, chatId)
	if conn.User.Id == authorId {
		return nil
	}
	conn.In <- e.LiveUpdate{
		Event:    event,
		ChatId:   chatId,
		UserId:   subjectId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[expelU]",
	}
	return nil
}

func chatDelete(conn *model.Conn, event e.UpdateType, authorId uint, chatID int, subjectId uint) error {
	log.Printf("∞----> chatDelete TRACE author[%d] deleted chat[%d] for subject[%d], target[%d]\n",
		authorId, chatID, subjectId, conn.User.Id)
	if subjectId != 0 && conn.User.Id != subjectId {
		return fmt.Errorf("chat_delete expect target and subject to be the same")
	}
	conn.In <- e.LiveUpdate{
		Event:    event,
		ChatId:   chatID,
		UserId:   conn.User.Id,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[deletedC]",
	}
	return nil
}

func chatClose(conn *model.Conn, event e.UpdateType, authorId uint, chatId int, subjectId uint) error {
	log.Printf("∞----> chatClose TRACE user[%d] closed chat[%d] for target[%d], subject[%d]\n",
		authorId, chatId, conn.User.Id, subjectId)
	if subjectId != 0 && conn.User.Id != subjectId {
		return fmt.Errorf("chat_close expect target and subject to be the same")
	}
	welcome := template.WelcomeTemplate{ActiveUser: conn.User.Name}
	data, err := welcome.HTML()
	if err != nil {
		return err
	}
	conn.In <- e.LiveUpdate{
		Event:    event,
		ChatId:   chatId,
		UserId:   conn.User.Id,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}
