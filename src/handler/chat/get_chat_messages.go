package handler

import (
	"fmt"
	"log"
	"neon-chat/src/app"
	"neon-chat/src/convert"
	"neon-chat/src/db"

	"github.com/jmoiron/sqlx"
)

func GetChatMessages(dbConn sqlx.Ext, chatId uint) ([]*app.Message, error) {
	chatUserIds, err := db.GetChatUserIds(dbConn, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] user ids, %s", chatId, err.Error())
	}
	dbUsers, err := db.GetUsers(dbConn, chatUserIds)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] users, %s", chatId, err.Error())
	}
	dbAvatars, err := db.GetAvatars(dbConn, chatUserIds)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] user avatars, %s", chatId, err.Error())
	}
	avatarByUserId := make(map[uint]*db.Avatar)
	for _, dbUser := range dbAvatars {
		avatarByUserId[dbUser.UserId] = dbUser
	}
	appUsers := make(map[uint]*app.User)
	for _, dbUser := range dbUsers {
		appUsers[dbUser.Id] = convert.UserDBToApp(&dbUser, avatarByUserId[dbUser.Id])
	}
	// TODO offset := 0 means no offset, ie get entire chat history
	dbMsgs, err := db.GetMessages(dbConn, chatId, 0)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] messages, %s", chatId, err.Error())
	}
	//
	appMsgs := make([]*app.Message, 0)
	appMsgIdMap := make(map[uint]*app.Message, 0)
	msgIds := make([]uint, len(dbMsgs))
	for _, dbMsg := range dbMsgs {
		author, ok := appUsers[dbMsg.AuthorId]
		if !ok {
			log.Printf("ERROR GetChatMessages author[%d] of message[%d] is not mapped\n", dbMsg.AuthorId, dbMsg.Id)
			continue
		}
		// ignore quote here
		appMsg := convert.MessageDBToApp(&dbMsg, author, nil)
		// sort the data on the way
		appMsgs = append(appMsgs, &appMsg)
		msgIds = append(msgIds, dbMsg.Id)
		appMsgIdMap[dbMsg.Id] = &appMsg
	}
	//
	dbQuotes, err := db.GetQuotes(dbConn, msgIds)
	if err != nil {
		return nil, fmt.Errorf("failed getting chat[%d] quotes, %s", chatId, err.Error())
	}
	//
	for _, dbQuote := range dbQuotes {
		appMsg, ok1 := appMsgIdMap[dbQuote.MsgId]
		appQuote, ok2 := appMsgIdMap[dbQuote.QuoteId]
		if ok1 && ok2 {
			appMsg.Quote = appQuote
		}
	}

	return appMsgs, nil
}
