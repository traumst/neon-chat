package pub

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"neon-chat/src/app"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/handler/priv"
)

func GetChats(dbConn sqlx.Ext, userId uint) ([]*app.Chat, error) {
	dbUserChats, err := db.GetUserChats(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user chats: %s", err.Error())
	}
	userChats := make([]*app.Chat, 0)
	for _, dbChat := range dbUserChats {
		appChatOwner, _ := priv.GetUser(dbConn, dbChat.OwnerId)
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
