package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/model/template"
	"neon-chat/src/state"
)

// TODO consider adding quote thread with depth limit
func GetQuote(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
	msgId uint,
) (string, error) {
	log.Printf("HandleGetQuote TRACE quoting message[%d] of chat[%d]\n", msgId, chatId)
	canChat, err := db.UsersCanChat(dbConn.Conn, chatId, user.Id)
	if err != nil {
		log.Printf("HandleGetQuote ERROR checking whether user[%d] can chat[%d], %s\n", user.Id, chatId, err)
		return "", fmt.Errorf("failed to check whether user can chat: %s", err.Error())
	} else if !canChat {
		log.Printf("HandleGetQuote ERROR user[%d] is not in chat[%d]\n", user.Id, chatId)
		return "", fmt.Errorf("user is not in chat")
	}
	//
	dbOwner, err := db.GetOwner(dbConn.Conn, chatId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting owner for chat[%d], %s\n", chatId, err.Error())
		return "", fmt.Errorf("failed to get chat owner from db, %s", err.Error())
	}
	appOwner := convert.UserDBToApp(dbOwner, nil)
	//
	dbMsg, err := db.GetMessage(dbConn.Conn, msgId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting quote[%d] from db, %s\n", msgId, err)
		return "", fmt.Errorf("failed to get message from db: %s", err.Error())
	}
	appQuote := convert.MessageDBToQuoteApp(dbMsg, user)
	//
	dbAvatar, err := db.GetAvatar(dbConn.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
		return "nil", fmt.Errorf("failed to get author avatar from db: %s", err.Error())
	}
	//
	dbAuthor, err := db.GetUser(dbConn.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting author[%d] from db, %s\n", dbMsg.AuthorId, err)
		return "", fmt.Errorf("failed to get author from db: %s", err.Error())
	}
	appQuote.Author = convert.UserDBToApp(dbAuthor, dbAvatar)
	//
	msgTmpl, err := appQuote.Template(user, appOwner)
	if err != nil {
		log.Printf("HandleGetQuote ERROR generating message template, %s\n", err)
		return "nil", fmt.Errorf("failed to generate message template: %s", err.Error())
	}
	quoteTmpl := &template.QuoteTemplate{
		IntermediateId: msgTmpl.IntermediateId,
		ChatId:         msgTmpl.ChatId,
		MsgId:          msgTmpl.MsgId,
		AuthorId:       msgTmpl.AuthorId,
		AuthorName:     msgTmpl.AuthorName,
		AuthorAvatar:   msgTmpl.AuthorAvatar,
		Text:           msgTmpl.Text,
		TextIntro:      msgTmpl.TextIntro,
	}
	return quoteTmpl.HTML()
}
