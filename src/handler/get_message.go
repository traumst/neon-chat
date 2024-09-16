package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	i "neon-chat/src/interface"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func GetMessage(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
	msgId uint,
) (i.Renderable, error) {
	log.Printf("HandleGetMessage TRACE getting message[%d] from chat[%d]\n", msgId, chatId)
	canChat, err := db.UsersCanChat(dbConn.Conn, chatId, user.Id)
	if err != nil {
		log.Printf("HandleGetMessage ERROR checking whether user[%d] can chat[%d], %s\n", user.Id, chatId, err)
		return nil, fmt.Errorf("failed to check whether user can chat: %s", err.Error())
	} else if !canChat {
		log.Printf("HandleGetMessage ERROR user[%d] is not in chat[%d]\n", user.Id, chatId)
		return nil, fmt.Errorf("user is not in chat")
	}
	dbOwner, err := db.GetOwner(dbConn.Conn, chatId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting owner for chat[%d], %s\n", chatId, err.Error())
		return nil, fmt.Errorf("failed to get message from db, %s", err.Error())
	}
	dbMsg, err := db.GetMessage(dbConn.Conn, msgId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting message[%d] from db, %s\n", msgId, err)
		return nil, fmt.Errorf("failed to get message from db: %s", err.Error())
	}
	dbAvatar, err := db.GetAvatar(dbConn.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
		return nil, fmt.Errorf("failed to get author avatar from db: %s", err.Error())
	}
	var appQuote *app.Message
	if dbQuote, _ := db.GetQuote(dbConn.Conn, msgId); dbQuote != nil {
		dbMsg, err := db.GetMessage(dbConn.Conn, dbQuote.QuoteId)
		if err != nil {
			log.Printf("HandleGetMessage warn getting quote message[%d] from db, %s\n", dbQuote.QuoteId, err)
			//return nil, fmt.Errorf("failed to get quote message[%d] from db: %s", dbQuote.QuoteId, err.Error())
		}
		dbAvatar, err := db.GetAvatar(dbConn.Conn, dbMsg.AuthorId)
		if err != nil {
			log.Printf("HandleGetMessage ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
			return nil, fmt.Errorf("failed to get author avatar from db: %s", err.Error())
		}
		user.Avatar = convert.AvatarDBToApp(dbAvatar)
		quoteMsg := convert.MessageDBToApp(dbMsg, user, nil)
		appQuote = &quoteMsg
	}
	appOwner := convert.UserDBToApp(dbOwner, nil)
	appMsg := convert.MessageDBToApp(dbMsg, user, appQuote)
	appMsg.Author.Avatar = convert.AvatarDBToApp(dbAvatar)
	tmplMsg, err := appMsg.Template(user, appOwner, appQuote)
	if err != nil {
		log.Printf("HandleGetMessage ERROR generating message template, %s\n", err)
		return nil, fmt.Errorf("failed to generate message template: %s", err.Error())
	}
	return &tmplMsg, nil
}
