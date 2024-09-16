package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func DeleteMessage(
	state *state.State,
	dbConn *db.DBConn,
	chatId uint,
	user *app.User,
	msgId uint,
) (*app.Chat, *app.Message, error) {
	log.Printf("HandleMessageDelete TRACE removing msg[%d] from chat[%d] for user[%d]\n", msgId, chatId, user.Id)
	dbMsg, err := db.GetMessage(dbConn.Tx, msgId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get message[%d] from db: %s", msgId, err.Error())
	}
	dbChat, err := db.GetChat(dbConn.Tx, chatId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chat[%d] from db: %s", chatId, err.Error())
	}
	if dbMsg.AuthorId != user.Id && dbChat.OwnerId != user.Id {
		return nil, nil, fmt.Errorf("user[%d] is not allowed to delete message[%d] in chat[%d]", user.Id, msgId, chatId)
	}
	err = db.DeleteMessage(dbConn.Tx, msgId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to remove message[%d] from chat[%d] in db, %s", msgId, chatId, err.Error())
	}
	dbSpecialUsers, err := db.GetUsers(dbConn.Tx, []uint{dbMsg.AuthorId, dbChat.OwnerId})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chat[%d] owner[%d] from db: %s", dbChat.Id, dbChat.OwnerId, err.Error())
	}
	var appChatOwner *app.User
	var appMsgAuthor *app.User
	for _, dbSpecialUser := range dbSpecialUsers {
		if dbSpecialUser.Id == dbChat.OwnerId {
			appChatOwner = convert.UserDBToApp(&dbSpecialUser, nil)
		}
		if dbSpecialUser.Id == dbMsg.AuthorId {
			appMsgAuthor = convert.UserDBToApp(&dbSpecialUser, nil)
		}
	}
	if appChatOwner == nil {
		return nil, nil, fmt.Errorf("chat[%d] owner[%d] not found", dbChat.Id, dbChat.OwnerId)
	}
	if appMsgAuthor == nil {
		return nil, nil, fmt.Errorf("message[%d] author[%d] not found", dbMsg.Id, dbMsg.AuthorId)
	}
	appChat := convert.ChatDBToApp(dbChat, &db.User{
		Id:     appChatOwner.Id,
		Name:   appChatOwner.Name,
		Email:  appChatOwner.Email,
		Type:   string(appChatOwner.Type),
		Status: string(appChatOwner.Status),
		Salt:   appChatOwner.Salt,
	})
	if appChat == nil {
		return nil, nil, fmt.Errorf("cannot convert chat[%d] for app, owner[%v]", dbChat.Id, appChatOwner)
	}
	appMsg := convert.MessageDBToApp(dbMsg, appMsgAuthor, nil) // TODO bad user
	return appChat, &appMsg, err
}
