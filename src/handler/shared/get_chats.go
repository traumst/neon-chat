package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	a "neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func GetChats(dbConn sqlx.Ext, userId uint) ([]*a.Chat, error) {
	dbUserChats, err := db.GetUserChats(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user chats: %s", err.Error())
	}
	userChats := make([]*a.Chat, 0)
	for _, dbChat := range dbUserChats {
		appChatOwner, _ := GetUser(dbConn, dbChat.OwnerId)
		if appChatOwner == nil {
			log.Printf("GetChats WARN chat[%d] owner[%d] not found\n", dbChat.Id, dbChat.OwnerId)
			continue
		}
		chat := convert.ChatDBToApp(&dbChat, &db.User{
			Id:     appChatOwner.Id,
			Type:   string(appChatOwner.Type),
			Status: string(appChatOwner.Status),
		})
		if chat == nil {
			log.Printf("GetChats WARN chat[%d] failed to map from db to app\n", dbChat.Id)
			continue
		}
		userChats = append(userChats, chat)
	}
	return userChats, nil
}
