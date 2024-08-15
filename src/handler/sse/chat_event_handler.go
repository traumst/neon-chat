package sse

import (
	"fmt"
	"log"

	"prplchat/src/handler/state"
	"prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
)

func chatCreate(conn *state.Conn, targetChat *app.Chat) error {
	log.Printf("chatCreate TRACE chat[%d] created by user[%d]\n", targetChat.Id, targetChat.OwnerId)
	template := targetChat.Template(conn.User, conn.User, []*app.User{conn.User}, []*app.Message{})
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveEvent{
		Event:    event.ChatAdd,
		ChatId:   targetChat.Id,
		UserId:   conn.User.Id,
		MsgId:    0,
		AuthorId: conn.User.Id,
		Data:     data,
	}
	return nil
}

func chatInvite(conn *state.Conn, targetChat *app.Chat, authorId uint, subject *app.User) error {
	log.Printf("chatCreate TRACE author[%d] invited subject[%d] to chat[%d], target[%d]\n",
		authorId, subject.Id, targetChat.Id, conn.User.Id)
	if authorId != targetChat.OwnerId {
		return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
			authorId, targetChat.OwnerId, targetChat.Id)
	}
	if subject == nil {
		return fmt.Errorf("subjectUser is nil for chatCreate")
	}
	if authorId == conn.User.Id || conn.User.Id != subject.Id {
		return fmt.Errorf("chatCreate conn[%s] does not belong to user[%d]", conn.Origin, subject.Id)
	}
	template := targetChat.Template(subject, conn.User, []*app.User{subject}, nil) // TODO BAD?
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveEvent{
		Event:    event.ChatInvite,
		ChatId:   targetChat.Id,
		UserId:   subject.Id,
		MsgId:    0,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}

func chatExpel(conn *state.Conn, chatId uint, ownerId uint, authorId uint, subjectId uint) error {
	log.Printf("chatExpel TRACE to user[%d] about author[%d] dropped subject[%d] from chat[%d]\n",
		conn.User.Id, authorId, subjectId, chatId)
	if authorId != ownerId && authorId != subjectId {
		return fmt.Errorf("author[%d] is not allowed to expel user[%d] from chat[%d]",
			authorId, subjectId, chatId)
	}
	conn.In <- event.LiveEvent{
		Event:    event.ChatExpel,
		ChatId:   chatId,
		UserId:   subjectId,
		MsgId:    0,
		AuthorId: authorId,
		Data:     "[expelU]",
	}
	return nil
}

func chatLeave(conn *state.Conn, chatId uint, ownerId uint, authorId uint, subjectId uint) error {
	log.Printf("chatLeave TRACE to user[%d] about author[%d] dropped subject[%d] from chat[%d]\n",
		conn.User.Id, authorId, subjectId, chatId)
	if authorId == ownerId || authorId != subjectId {
		return fmt.Errorf("author[%d] is not allowed to leave chat[%d] for user[%d]",
			authorId, chatId, subjectId)
	}
	conn.In <- event.LiveEvent{
		Event:    event.ChatLeave,
		ChatId:   chatId,
		UserId:   subjectId,
		MsgId:    0,
		AuthorId: authorId,
		Data:     "[leftU]",
	}
	return nil
}

func chatDelete(conn *state.Conn, chatId uint, ownerId uint, authorId uint, targetId uint) error {
	log.Printf("chatDelete TRACE deleted chat[%d] for subject[%d], target[%d]\n",
		chatId, targetId, conn.User.Id)
	if authorId != ownerId && authorId != targetId {
		return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
			authorId, ownerId, chatId)
	}
	if targetId != 0 && conn.User.Id != targetId {
		return fmt.Errorf("chatDelete conn[%s] does not belong to user[%d]", conn.Origin, targetId)
	}
	conn.In <- event.LiveEvent{
		Event:    event.ChatDrop,
		ChatId:   chatId,
		UserId:   targetId,
		MsgId:    0,
		AuthorId: authorId,
		Data:     "[deletedC]",
	}
	return nil
}

func chatClose(conn *state.Conn, chatId uint, ownerId uint, authorId uint, target *app.User) error {
	log.Printf("chatClose TRACE user[%d] closed chat[%d] for subject[%d]\n", authorId, chatId, target.Id)
	if authorId != ownerId && authorId != target.Id {
		return fmt.Errorf("author[%d] is not allowed to close chat[%d] for user[%d]",
			authorId, chatId, target.Id)
	}
	if target.Id != 0 && conn.User.Id != target.Id {
		return fmt.Errorf("chatClose conn[%s] belongs to other user[%d]", conn.Origin, conn.User.Id)
	}
	welcome := template.WelcomeTemplate{User: template.UserTemplate{}}
	data, err := welcome.HTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveEvent{
		Event:    event.ChatClose,
		ChatId:   chatId,
		UserId:   target.Id,
		MsgId:    0,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}
