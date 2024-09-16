package chat

import (
	"fmt"
	"log"

	"neon-chat/src/convert"
	"neon-chat/src/db"
	u "neon-chat/src/handler/user"
	"neon-chat/src/model/app"
	"neon-chat/src/model/template"
	"neon-chat/src/state"
)

func AddChat(state *state.State, dbConn *db.DBConn, user *app.User, chatName string) (*template.ChatTemplate, error) {
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
	appChatUsers, err := u.GetChatUsers(dbConn.Tx, dbChat.Id)
	if err != nil {
		log.Printf("HandleChatAdd ERROR getting chat[%d] users: %s", dbChat.Id, err)
	}
	appChatMsgs := make([]*app.Message, 0)
	tmpl := appChat.Template(user, user, appChatUsers, appChatMsgs)
	return &tmpl, nil
}
