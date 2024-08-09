package handler

import (
	"fmt"
	"log"
	"prplchat/src/convert"
	d "prplchat/src/db"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
)

func HandleMessageAdd(
	state *state.State,
	db *d.DBConn,
	chatId uint,
	author *a.User,
	msg string,
) (*a.Message, error) {
	log.Printf("HandleMessageAdd TRACE opening current chat for user[%d]\n", author.Id)
	canChat, err := db.UserCanChat(chatId, author.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to check user[%d] can chat[%d]: %s", author.Id, chatId, err.Error())
	}
	if !canChat {
		return nil, fmt.Errorf("user is not in chat")
	}
	dbChat, err := db.GetChat(chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] from db: %s", chatId, err.Error())
	}
	dbMsg, err := db.AddMessage(&d.Message{
		Id:       0,
		ChatId:   chatId,
		AuthorId: author.Id,
		Text:     msg, // TODO: sanitize
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	appChat := convert.ChatDBToApp(dbChat)
	appMsg := convert.MessageDBToApp(dbMsg, author)
	err = sse.DistributeMsg(state, db, appChat, author.Id, &appMsg, event.MessageAdd)
	if err != nil {
		log.Printf("HandleMessageAdd ERROR distributing msg update, %s\n", err)
	}
	return &appMsg, nil
}

func HandleMessageDelete(
	state *state.State,
	db *d.DBConn,
	chatId uint,
	user *a.User,
	msgId uint,
) (*a.Message, error) {
	log.Printf("HandleMessageDelete TRACE removing msg[%d] from chat[%d] for user[%d]\n", msgId, chatId, user.Id)
	dbMsg, err := db.GetMessage(msgId)
	if err != nil {
		return nil, fmt.Errorf("failed to get message[%d] from db: %s", msgId, err.Error())
	}
	dbChat, err := db.GetChat(chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] from db: %s", chatId, err.Error())
	}
	if dbMsg.AuthorId != user.Id && dbChat.OwnerId != user.Id {
		return nil, fmt.Errorf("user[%d] is not allowed to delete message[%d] in chat[%d]", user.Id, msgId, chatId)
	}
	err = db.DeleteMessage(msgId)
	if err != nil {
		return nil, fmt.Errorf("failed to remove message[%d] from chat[%d] in db, %s", msgId, chatId, err.Error())
	}
	appChat := convert.ChatDBToApp(dbChat)
	appMsg := convert.MessageDBToApp(dbMsg, &a.User{Id: dbMsg.AuthorId}) // TODO bad user
	err = sse.DistributeMsg(state, db, appChat, user.Id, &appMsg, event.MessageDrop)
	if err != nil {
		log.Printf("HandleMessageDelete ERROR distributing msg update, %s\n", err)
	}
	return &appMsg, err
}
