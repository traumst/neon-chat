package pub

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

// TODO consider quote thread with depth limit
func GetQuote(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatId uint,
	msgId uint,
) (*app.Quote, error) {
	log.Printf("HandleGetQuote TRACE quoting message[%d] of chat[%d]\n", msgId, chatId)
	canChat, err := db.UsersCanChat(dbConn.Conn, chatId, user.Id)
	if err != nil {
		log.Printf("HandleGetQuote ERROR checking whether user[%d] can chat[%d], %s\n", user.Id, chatId, err)
		return nil, fmt.Errorf("failed to check whether user can chat: %s", err.Error())
	} else if !canChat {
		log.Printf("HandleGetQuote ERROR user[%d] is not in chat[%d]\n", user.Id, chatId)
		return nil, fmt.Errorf("user is not in chat")
	}
	//
	dbMsg, err := db.GetMessage(dbConn.Conn, msgId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting quote[%d] from db, %s\n", msgId, err)
		return nil, fmt.Errorf("failed to get message from db: %s", err.Error())
	}
	appQuote := convert.MessageDBToQuoteApp(dbMsg, user)
	//
	dbAvatar, err := db.GetAvatar(dbConn.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
		return nil, fmt.Errorf("failed to get author avatar from db: %s", err.Error())
	}
	//
	dbAuthor, err := db.GetUser(dbConn.Conn, dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetQuote ERROR getting author[%d] from db, %s\n", dbMsg.AuthorId, err)
		return nil, fmt.Errorf("failed to get author from db: %s", err.Error())
	}
	appQuote.Author = convert.UserDBToApp(dbAuthor, dbAvatar)
	return &appQuote, nil
}
