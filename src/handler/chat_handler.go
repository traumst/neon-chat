package handler

import (
	"fmt"
	"log"

	"prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
)

func chatCreate(conn *Conn, targetChat *app.Chat, author *app.User) error {
	log.Printf("chatCreate TRACE author[%d] created chat[%d]\n",
		author.Id, targetChat.Id)
	if author.Id != targetChat.Owner.Id {
		return fmt.Errorf("author[%d] is not owner[%d] of chat[%d]",
			author.Id, targetChat.Owner.Id, targetChat.Id)
	}
	if author.Id != conn.User.Id || conn.User.Id != author.Id {
		return fmt.Errorf("chatCreate conn[%s] does not belong to user[%d]", conn.Origin, author.Id)
	}
	template := targetChat.Template(author, conn.User)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveEvent{
		Event:    event.ChatAdd,
		ChatId:   targetChat.Id,
		UserId:   author.Id,
		MsgId:    -1,
		AuthorId: author.Id,
		Data:     data,
	}
	return nil
}

func chatInvite(conn *Conn, targetChat *app.Chat, authorId uint, subject *app.User) error {
	log.Printf("chatCreate TRACE author[%d] invited subject[%d] to chat[%d], target[%d]\n",
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
	template := targetChat.Template(subject, conn.User)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveEvent{
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
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[expelU]",
	}
	return nil
}

func chatLeave(conn *Conn, chatId int, ownerId uint, authorId uint, subjectId uint) error {
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
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[leftU]",
	}
	return nil
}

func chatDelete(conn *Conn, chatId int, ownerId uint, authorId uint, targetId uint) error {
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
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[deletedC]",
	}
	return nil
}

func chatClose(conn *Conn, chatId int, ownerId uint, authorId uint, target *app.User) error {
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
		MsgId:    -1,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}
