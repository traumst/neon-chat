package handler

import (
	"fmt"
	d "prplchat/src/db"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
)

func HandleChatAdd(app *state.State, db *d.DBConn, user *a.User, chatName string) (*a.Chat, error) {
	chat, err := db.AddChat(&d.Chat{Title: chatName, OwnerId: user.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to add chat to db: %s", err)
	}
	err = app.AddChat(chat.Id, chat.Title, user)
	if err != nil {
		return nil, fmt.Errorf("failed to add chat to app")
	}
	openChat, err := app.OpenChat(user.Id, chat.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to open new chat: %s", err)
	}
	return openChat, nil
}

func HandleChatDelete(app *state.State, db *d.DBConn, userId uint, chat *a.Chat) error {
	err := db.DeleteChat(chat.Id)
	if err != nil {
		return fmt.Errorf("error deleting chat in db: %s", err.Error())
	}
	err = app.DeleteChat(userId, chat)
	if err != nil {
		return fmt.Errorf("error deleting chat in app: %s", err.Error())
	}
	return nil
}
