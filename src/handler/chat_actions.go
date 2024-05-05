package handler

import (
	"fmt"
	"log"

	"go.chat/src/model/app"
	"go.chat/src/model/event"
	"go.chat/src/model/template"
)

func chatCreate(conn *Conn, targetChat *app.Chat, author *app.User) error {
	log.Printf("∞----> chatCreate TRACE author[%d] created chat[%d]\n",
		author.Id, targetChat.Id)
	if author.Id != targetChat.Owner.Id {
		return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
			author.Id, targetChat.Owner.Id, targetChat.Id)
	}
	if author.Id != conn.User.Id || conn.User.Id != author.Id {
		return fmt.Errorf("chatCreate conn[%s] does not belong to user[%d]", conn.Origin, author.Id)
	}
	template := targetChat.Template(author)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveUpdate{
		Event:    event.ChatCreated,
		ChatId:   targetChat.Id,
		UserId:   author.Id,
		MsgId:    -1,
		AuthorId: author.Id,
		Data:     data,
	}
	return nil
}

func chatInvite(conn *Conn, targetChat *app.Chat, authorId uint, subject *app.User) error {
	log.Printf("∞----> chatCreate TRACE author[%d] invited subject[%d] to chat[%d], target[%d]\n",
		authorId, subject.Id, targetChat.Id, conn.User.Id)
	if authorId != targetChat.Owner.Id {
		return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
			authorId, targetChat.Owner.Id, targetChat.Id)
	}
	if subject == nil {
		return fmt.Errorf("subjectUser is nil for chatCreate")
	}
	if authorId == conn.User.Id || conn.User.Id != subject.Id {
		return fmt.Errorf("chatCreate conn[%s] does not belong to user[%d]", conn.Origin, subject.Id)
	}
	template := targetChat.Template(subject)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveUpdate{
		Event:    event.ChatInvite,
		ChatId:   targetChat.Id,
		UserId:   subject.Id,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}

func chatExpel(conn *Conn, chatId int, ownerId uint, authorId uint, subjectId uint) error {
	log.Printf("∞----> chatExpel TRACE to user[%d] about author[%d] dropped subject[%d] from chat[%d]\n",
		conn.User.Id, authorId, subjectId, chatId)
	if authorId != ownerId && authorId != subjectId {
		return fmt.Errorf("author[%d] is not allowed to expel user[%d] from chat[%d]",
			authorId, subjectId, chatId)
	}
	conn.In <- event.LiveUpdate{
		Event:    event.ChatExpel,
		ChatId:   chatId,
		UserId:   subjectId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[expelU]",
	}
	return nil
}

func chatLeave(conn *Conn, chatId int, ownerId uint, authorId uint, subjectId uint) error {
	log.Printf("∞----> chatLeave TRACE to user[%d] about author[%d] dropped subject[%d] from chat[%d]\n",
		conn.User.Id, authorId, subjectId, chatId)
	if authorId == ownerId || authorId != subjectId {
		return fmt.Errorf("author[%d] is not allowed to leave chat[%d] for user[%d]",
			authorId, chatId, subjectId)
	}
	conn.In <- event.LiveUpdate{
		Event:    event.ChatLeave,
		ChatId:   chatId,
		UserId:   subjectId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[leftU]",
	}
	return nil
}

func chatDelete(conn *Conn, chatId int, ownerId uint, authorId uint, targetId uint) error {
	log.Printf("∞----> chatDelete TRACE deleted chat[%d] for subject[%d], target[%d]\n",
		chatId, targetId, conn.User.Id)
	if authorId != ownerId && authorId != targetId {
		return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
			authorId, ownerId, chatId)
	}
	if targetId != 0 && conn.User.Id != targetId {
		return fmt.Errorf("chatDelete conn[%s] does not belong to user[%d]", conn.Origin, targetId)
	}
	conn.In <- event.LiveUpdate{
		Event:    event.ChatDeleted,
		ChatId:   chatId,
		UserId:   targetId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[deletedC]",
	}
	return nil
}

func chatClose(conn *Conn, chatId int, ownerId uint, authorId uint, targetId uint) error {
	log.Printf("∞----> chatClose TRACE user[%d] closed chat[%d] for subject[%d]\n", authorId, chatId, targetId)
	if authorId != ownerId && authorId != targetId {
		return fmt.Errorf("author[%d] is not allowed to close chat[%d] for user[%d]",
			authorId, chatId, targetId)
	}
	if targetId != 0 && conn.User.Id != targetId {
		return fmt.Errorf("chatClose conn[%s] belongs to other user[%d]", conn.Origin, conn.User.Id)
	}
	welcome := template.WelcomeTemplate{ActiveUser: conn.User.Name}
	data, err := welcome.HTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveUpdate{
		Event:    event.ChatClose,
		ChatId:   chatId,
		UserId:   targetId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}
