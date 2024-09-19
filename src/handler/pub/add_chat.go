package pub

import (
	"fmt"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/state"
)

func AddChat(
	state *state.State,
	dbConn *db.DBConn,
	user *app.User,
	chatName string,
) (*app.Chat, error) {
	dbChat, err := db.AddChat(dbConn.Tx, &db.Chat{Title: chatName, OwnerId: user.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to add chat[%s] to db: %s", chatName, err)
	}
	err = db.AddChatUser(dbConn.Tx, dbChat.Id, user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to add owner[%d] to chat[%d] in db: %s", user.Id, dbChat.Id, err)
	}
	err = state.OpenChat(user.Id, dbChat.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to open new chat: %s", err)
	}
	openChat, err := db.GetChat(dbConn.Tx, dbChat.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat[%d] from db: %s", dbChat.Id, err)
	}
	appChat := convert.ChatDBToApp(openChat, &db.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	})
	return appChat, nil
}
