package handler

import (
	"fmt"
	"log"
	d "prplchat/src/db"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
)

func HandleMessageAdd(
	app *state.State,
	db *d.DBConn,
	chatId uint,
	author *a.User,
	msg string,
) (*a.Message, error) {
	log.Printf("HandleMessageAdd TRACE opening current chat for user[%d]\n", author.Id)
	if canChat, _ := db.UserCanChat(chatId, author.Id); !canChat {
		return nil, fmt.Errorf("user is not in chat")
	}
	dbMsg, err := addMsgIntoDB(db, chatId, author, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to add message to db: %s", err.Error())
	}
	appChat, err := app.GetChat(dbMsg.AuthorId, dbMsg.ChatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat from app: %s", err.Error())
	}
	appMsg, err := addMsgIntoApp(app, dbMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to add message to app: %s", err.Error())
	}
	err = sse.DistributeMsg(app, appChat, author.Id, appMsg, event.MessageAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to distribute new message, %s", err.Error())
	}
	return appMsg, nil
}

func addMsgIntoDB(
	db *d.DBConn,
	chatId uint,
	author *a.User,
	msg string,
) (*d.Message, error) {
	log.Printf("addMsgIntoApp TRACE storing message for user[%d] in chat[%d]\n", author.Id, chatId)
	dbChat, err := db.GetChat(chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %s", err.Error())
	}
	dbMsg, err := db.AddMessage(&d.Message{
		Id:       0,
		ChatId:   chatId,
		OwnerId:  dbChat.OwnerId,
		AuthorId: author.Id,
		Text:     msg, // TODO: sanitize
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	return dbMsg, nil
}

func addMsgIntoApp(
	app *state.State,
	dbMsg *d.Message,
) (*a.Message, error) {
	log.Printf("addMsgIntoApp TRACE storing message for user[%d] in chat[%d]\n", dbMsg.AuthorId, dbMsg.ChatId)
	appChat, err := app.GetChat(dbMsg.AuthorId, dbMsg.ChatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat from app: %s", err.Error())
	}
	newMsg := MessageDBToApp(dbMsg)
	appMsg, err := appChat.AddMessage(dbMsg.AuthorId, newMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to store message: %s", err.Error())
	}
	return appMsg, nil
}
