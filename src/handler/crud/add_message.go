package crud

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func AddMessage(
	state *state.State,
	dbConn *db.DBConn,
	chatId uint,
	author *app.User,
	msg string,
	quoteId uint,
) (*app.Chat, *app.Message, error) {
	log.Printf("HandleMessageAdd TRACE opening current chat for user[%d]\n", author.Id)
	canChat, err := db.UsersCanChat(dbConn.Tx, chatId, author.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check user[%d] can chat[%d]: %s", author.Id, chatId, err.Error())
	}
	if !canChat {
		return nil, nil, fmt.Errorf("user is not in chat")
	}
	dbChat, err := db.GetChat(dbConn.Tx, chatId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chat[%d] from db: %s", chatId, err.Error())
	}
	dbMsg, err := db.AddMessage(dbConn.Tx, &db.Message{
		Id:       0,
		ChatId:   chatId,
		AuthorId: author.Id,
		Text:     msg,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	// quoteId 0 means message has no quote attached
	var appQuote *app.Message
	if quoteId != 0 {
		quote, err := db.AddQuote(dbConn.Tx, &db.Quote{
			MsgId:   dbMsg.Id,
			QuoteId: quoteId,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to add quote[%d] to message[%d]: %s", quoteId, dbMsg.Id, err.Error())
		}
		dbQuote, err := db.GetMessage(dbConn.Tx, quote.QuoteId)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get quote[%d] from db: %s", quoteId, err.Error())
		}
		quoteAuthor, err := GetUser(dbConn.Tx, dbQuote.AuthorId)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get quote[%d] author[%d] avatar from db: %s", quoteId, dbQuote.AuthorId, err.Error())
		}
		tmp := convert.MessageDBToApp(dbQuote, quoteAuthor, nil)
		appQuote = &tmp
	}
	dbOwner, err := db.GetUser(dbConn.Tx, dbChat.OwnerId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chat[%d] owner[%d] from db: %s", chatId, dbChat.OwnerId, err.Error())
	}
	appChat := convert.ChatDBToApp(dbChat, dbOwner)
	appMsg := convert.MessageDBToApp(dbMsg, author, appQuote)
	return appChat, &appMsg, nil
}
