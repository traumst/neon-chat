package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler/sse"
	"neon-chat/src/handler/state"
	i "neon-chat/src/interface"
	a "neon-chat/src/model/app"
	"neon-chat/src/model/event"
)

func HandleGetMessage(
	state *state.State,
	db *d.DBConn,
	user *a.User,
	chatId uint,
	msgId uint,
) (i.Renderable, error) {
	log.Printf("HandleGetMessage TRACE opening current chat for user[%d]\n", user.Id)
	canChat, err := db.UserCanChat(chatId, user.Id)
	if err != nil {
		log.Printf("HandleGetMessage ERROR checking user[%d] can chat[%d], %s\n", user.Id, chatId, err)
		return nil, fmt.Errorf("failed to check whether user can chat: %s", err.Error())
	} else if !canChat {
		log.Printf("HandleGetMessage ERROR user[%d] is not in chat[%d]\n", user.Id, chatId)
		return nil, fmt.Errorf("user is not in chat")
	}

	dbOwner, err := db.GetOwner(msgId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting ownerfor chat[%d], %s\n", chatId, err.Error())
		return nil, fmt.Errorf("failed to get message from db, %s", err.Error())
	}
	dbMsg, err := db.GetMessage(msgId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting message[%d] from db, %s\n", msgId, err)
		return nil, fmt.Errorf("failed to get message from db: %s", err.Error())
	}
	dbAvatar, err := db.GetAvatar(dbMsg.AuthorId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting author[%d] avatar from db, %s\n", dbMsg.AuthorId, err)
		return nil, fmt.Errorf("failed to get author avatar from db: %s", err.Error())
	}

	appOwner := convert.UserDBToApp(dbOwner)
	appMsg := convert.MessageDBToApp(dbMsg, user)
	appAvatar := convert.AvatarDBToApp(dbAvatar)
	tmplMsg, err := appMsg.Template(user, appOwner, appAvatar)
	if err != nil {
		log.Printf("HandleGetMessage ERROR generating message template, %s\n", err)
		return nil, fmt.Errorf("failed to generate message template: %s", err.Error())
	}

	return tmplMsg, nil
}

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
	//userInput := tokenizer.ParsedUserInput(msg)
	dbMsg, err := db.AddMessage(&d.Message{
		Id:       0,
		ChatId:   chatId,
		AuthorId: author.Id,
		Text:     msg, // TODO: sanitize
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	dbOwner, err := db.GetUser(dbChat.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] owner[%d] from db: %s", chatId, dbChat.OwnerId, err.Error())
	}
	appChat := convert.ChatDBToApp(dbChat, convert.UserDBToApp(dbOwner))
	appMsg := convert.MessageDBToApp(dbMsg, author)
	err = sse.DistributeMsg(state, db, appChat, &appMsg, event.MessageAdd)
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
	dbSpecialUsers, err := db.GetUsers([]uint{dbMsg.AuthorId, dbChat.OwnerId})
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] owner[%d] from db: %s", dbChat.Id, dbChat.OwnerId, err.Error())
	}
	var appChatOwner *a.User
	var appMsgAuthor *a.User
	for _, dbSpecialUser := range dbSpecialUsers {
		if dbSpecialUser.Id == dbChat.OwnerId {
			appChatOwner = convert.UserDBToApp(&dbSpecialUser)
		}
		if dbSpecialUser.Id == dbMsg.AuthorId {
			appMsgAuthor = convert.UserDBToApp(&dbSpecialUser)
		}
	}
	if appChatOwner == nil {
		return nil, fmt.Errorf("chat[%d] owner[%d] not found", dbChat.Id, dbChat.OwnerId)
	}
	if appMsgAuthor == nil {
		return nil, fmt.Errorf("message[%d] author[%d] not found", dbMsg.Id, dbMsg.AuthorId)
	}
	appChat := convert.ChatDBToApp(dbChat, appChatOwner)
	if appChat == nil {
		return nil, fmt.Errorf("cannot convert chat[%d] for app, owner[%v]", dbChat.Id, appChatOwner)
	}
	appMsg := convert.MessageDBToApp(dbMsg, appMsgAuthor) // TODO bad user
	err = sse.DistributeMsg(state, db, appChat, &appMsg, event.MessageDrop)
	if err != nil {
		log.Printf("HandleMessageDelete ERROR distributing msg update, %s\n", err)
	}
	return &appMsg, err
}

func GetChatMessages(db *d.DBConn, chatId uint) ([]*a.Message, error) {
	// TODO offset := 0 means no offset, ie get entire chat history
	dbMsgs, err := db.GetMessages(chatId, 0)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] messages, %s", chatId, err.Error())
	}
	msgAuthorIds := make(map[uint]bool)
	authorIds := make([]uint, 0)
	for _, dbMsg := range dbMsgs {
		if msgAuthorIds[dbMsg.AuthorId] {
			continue
		}
		msgAuthorIds[dbMsg.AuthorId] = true
		authorIds = append(authorIds, dbMsg.AuthorId)
	}
	dbUsers, err := db.GetUsers(authorIds)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] users[%v], %s", chatId, authorIds, err.Error())
	}
	authorAvatars := make(map[uint]*a.Avatar)
	dbAvatars, err := db.GetAvatars(authorIds)
	if err == nil {
		for _, dbAvatar := range dbAvatars {
			if authorAvatars[dbAvatar.UserId] != nil {
				continue
			}
			authorAvatars[dbAvatar.UserId] = convert.AvatarDBToApp(dbAvatar)
		}
	}
	msgAuthors := make(map[uint]*a.User)
	for _, dbUser := range dbUsers {
		msgAuthors[dbUser.Id] = convert.UserDBToApp(&dbUser)
	}
	appChatMsgs := make([]*a.Message, 0)
	for _, dbMsg := range dbMsgs {
		author, ok := msgAuthors[dbMsg.AuthorId]
		if !ok {
			log.Printf("GetChatMessages ERROR author[%d] of message[%d] is not mapped\n", dbMsg.AuthorId, dbMsg.Id)
			continue
		}
		if authorAvatars[author.Id] != nil {
			author.Avatar = authorAvatars[author.Id]
		}
		appMsg := convert.MessageDBToApp(&dbMsg, author)
		appChatMsgs = append(appChatMsgs, &appMsg)
	}
	return appChatMsgs, nil
}
