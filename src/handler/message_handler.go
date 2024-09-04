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

// TODO consider adding quote thread with depth limit
func HandleGetMessage(
	state *state.State,
	db *d.DBConn,
	user *a.User,
	chatId uint,
	msgId uint,
) (i.Renderable, error) {
	log.Printf("HandleGetMessage TRACE opening current chat for user[%d]\n", user.Id)
	canChat, err := db.UsersCanChat(chatId, user.Id)
	if err != nil {
		log.Printf("HandleGetMessage ERROR checking user[%d] can chat[%d], %s\n", user.Id, chatId, err)
		return nil, fmt.Errorf("failed to check whether user can chat: %s", err.Error())
	} else if !canChat {
		log.Printf("HandleGetMessage ERROR user[%d] is not in chat[%d]\n", user.Id, chatId)
		return nil, fmt.Errorf("user is not in chat")
	}

	dbOwner, err := db.GetOwner(chatId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting owner for chat[%d], %s\n", chatId, err.Error())
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
	dbQuote, err := db.GetQuote(msgId)
	if err != nil {
		log.Printf("HandleGetMessage ERROR getting quotes for message[%d], %s\n", msgId, err)
		return nil, fmt.Errorf("failed to get quotes for message: %s", err.Error())
	}
	var appQuote *a.Message
	if dbQuote != nil {
		dbMsg, err := db.GetMessage(dbQuote.QuoteId)
		if err != nil {
			log.Printf("HandleGetMessage ERROR getting quote message[%d] from db, %s\n", dbQuote.QuoteId, err)
			return nil, fmt.Errorf("failed to get quote message from db: %s", err.Error())
		}
		quoteMsg := convert.MessageDBToApp(dbMsg, user, nil)
		appQuote = &quoteMsg
	}

	appOwner := convert.UserDBToApp(dbOwner)
	appMsg := convert.MessageDBToApp(dbMsg, user, appQuote)
	appAvatar := convert.AvatarDBToApp(dbAvatar)
	tmplMsg, err := appMsg.Template(user, appOwner, appAvatar, appQuote)
	if err != nil {
		log.Printf("HandleGetMessage ERROR generating message template, %s\n", err)
		return nil, fmt.Errorf("failed to generate message template: %s", err.Error())
	}

	return &tmplMsg, nil
}

// TODO has partial success state - if fail of saving quote
func HandleMessageAdd(
	state *state.State,
	db *d.DBConn,
	chatId uint,
	author *a.User,
	msg string,
	quoteId uint,
) (*a.Message, error) {
	log.Printf("HandleMessageAdd TRACE opening current chat for user[%d]\n", author.Id)
	canChat, err := db.UsersCanChat(chatId, author.Id)
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
		Text:     msg,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add message to chat[%d]: %s", chatId, err.Error())
	}
	// quoteId 0 means message has no quote attached
	var appQuote *a.Message
	if quoteId != 0 {
		_, err = db.AddQuote(&d.Quote{
			MsgId:   dbMsg.Id,
			QuoteId: quoteId,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add quote[%d] to message[%d]: %s", quoteId, dbMsg.Id, err.Error())
		}
		dbQuote, err := db.GetMessage(quoteId)
		if err != nil {
			return nil, fmt.Errorf("failed to get quote[%d] from db: %s", quoteId, err.Error())
		}
		quote := convert.MessageDBToApp(dbQuote, author, nil)
		appQuote = &quote
	}
	dbOwner, err := db.GetUser(dbChat.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] owner[%d] from db: %s", chatId, dbChat.OwnerId, err.Error())
	}
	appChat := convert.ChatDBToApp(dbChat, convert.UserDBToApp(dbOwner))
	appMsg := convert.MessageDBToApp(dbMsg, author, appQuote)
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
	appMsg := convert.MessageDBToApp(dbMsg, appMsgAuthor, nil) // TODO bad user
	err = sse.DistributeMsg(state, db, appChat, &appMsg, event.MessageDrop)
	if err != nil {
		log.Printf("HandleMessageDelete ERROR distributing msg update, %s\n", err)
	}
	return &appMsg, err
}

func GetChatMessages(db *d.DBConn, chatId uint) ([]*a.Message, error) {
	dbUsers, err := db.GetChatUsers(chatId)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] users, %s", chatId, err.Error())
	}
	userIds := make([]uint, 0)
	appUsers := make(map[uint]*a.User)
	for _, dbUser := range dbUsers {
		appUsers[dbUser.Id] = convert.UserDBToApp(&dbUser)
		userIds = append(userIds, dbUser.Id)
	}
	appAvatars := make(map[uint]*a.Avatar)
	dbAvatars, err := db.GetAvatars(userIds)
	if err == nil {
		for _, dbAvatar := range dbAvatars {
			if appAvatars[dbAvatar.UserId] != nil {
				continue
			}
			appAvatars[dbAvatar.UserId] = convert.AvatarDBToApp(dbAvatar)
		}
	}
	// TODO offset := 0 means no offset, ie get entire chat history
	dbMsgs, err := db.GetMessages(chatId, 0)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] messages, %s", chatId, err.Error())
	}
	appMsgs := make([]*a.Message, 0)
	appMsgIdMap := make(map[uint]*a.Message, 0)
	msgIds := make([]uint, len(dbMsgs))
	for _, dbMsg := range dbMsgs {
		author, ok := appUsers[dbMsg.AuthorId]
		if !ok {
			log.Printf("GetChatMessages ERROR author[%d] of message[%d] is not mapped\n", dbMsg.AuthorId, dbMsg.Id)
			continue
		}
		if appAvatars[author.Id] != nil {
			author.Avatar = appAvatars[author.Id]
		}
		// ignore quote for now
		appMsg := convert.MessageDBToApp(&dbMsg, author, nil)
		// sort the data on the way
		appMsgs = append(appMsgs, &appMsg)
		msgIds = append(msgIds, dbMsg.Id)
		appMsgIdMap[dbMsg.Id] = &appMsg
	}
	dbQuotes, err := db.GetQuotes(msgIds)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] quotes, %s", chatId, err.Error())
	}
	for _, dbQuote := range dbQuotes {
		appMsg, ok1 := appMsgIdMap[dbQuote.MsgId]
		appQuote, ok2 := appMsgIdMap[dbQuote.QuoteId]
		if ok1 && ok2 {
			appMsg.Quote = appQuote
		}
	}

	return appMsgs, nil
}
