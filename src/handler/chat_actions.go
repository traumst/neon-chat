package handler

import (
	"fmt"
	"log"

	"go.chat/src/model/app"
	"go.chat/src/model/event"
	"go.chat/src/model/template"
)

func chatCreate(conn *Conn, evnt event.UpdateType, targetChat *app.Chat, author *app.User) error {
	log.Printf("∞----> chatCreate TRACE author[%d] created chat[%d]\n",
		author.Id, targetChat.Id)
	if author.Id != conn.User.Id || conn.User.Id != author.Id {
		return fmt.Errorf("chatCreate conn[%s] does not belong to user[%d]", conn.Origin, author.Id)
	}
	template := targetChat.Template(author)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveUpdate{
		Event:    evnt,
		ChatId:   targetChat.Id,
		UserId:   author.Id,
		MsgId:    -1,
		AuthorId: author.Id,
		Data:     data,
	}
	return nil
}

func chatInvite(conn *Conn, evnt event.UpdateType, targetChat *app.Chat, authorId uint, subject *app.User) error {
	log.Printf("∞----> chatCreate TRACE author[%d] invited subject[%d] to chat[%d], target[%d]\n",
		authorId, subject.Id, targetChat.Id, conn.User.Id)
	if authorId == conn.User.Id || conn.User.Id != subject.Id {
		return fmt.Errorf("chatCreate conn[%s] does not belong to user[%d]", conn.Origin, subject.Id)
	}
	template := targetChat.Template(subject)
	data, err := template.ShortHTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveUpdate{
		Event:    evnt,
		ChatId:   targetChat.Id,
		UserId:   subject.Id,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}

func chatExpel(conn *Conn, evnt event.UpdateType, chatId int, authorId uint, subjectId uint) error {
	log.Printf("∞----> chatExpel TRACE to user[%d] about author[%d] dropped subject[%d] from chat[%d]\n",
		conn.User.Id, authorId, subjectId, chatId)
	conn.In <- event.LiveUpdate{
		Event:    evnt,
		ChatId:   chatId,
		UserId:   subjectId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[expelU]",
	}
	return nil
}

func chatDelete(conn *Conn, evnt event.UpdateType, chatId int, authorId uint, targetId uint) error {
	log.Printf("∞----> chatDelete TRACE deleted chat[%d] for subject[%d], target[%d]\n",
		chatId, targetId, conn.User.Id)
	if targetId != 0 && conn.User.Id != targetId {
		return fmt.Errorf("chatDelete conn[%s] does not belong to user[%d]", conn.Origin, targetId)
	}
	conn.In <- event.LiveUpdate{
		Event:    evnt,
		ChatId:   chatId,
		UserId:   targetId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     "[deletedC]",
	}
	return nil
}

func chatClose(conn *Conn, evnt event.UpdateType, chatId int, authorId uint, targetId uint) error {
	log.Printf("∞----> chatClose TRACE user[%d] closed chat[%d] for subject[%d]\n", authorId, chatId, targetId)
	if targetId != 0 && conn.User.Id != targetId {
		return fmt.Errorf("chatClose conn[%s] belongs to other user[%d]", conn.Origin, conn.User.Id)
	}
	welcome := template.WelcomeTemplate{ActiveUser: conn.User.Name}
	data, err := welcome.HTML()
	if err != nil {
		return err
	}
	conn.In <- event.LiveUpdate{
		Event:    evnt,
		ChatId:   chatId,
		UserId:   targetId,
		MsgId:    -1,
		AuthorId: authorId,
		Data:     data,
	}
	return nil
}
