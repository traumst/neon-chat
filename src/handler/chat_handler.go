package handler

import (
	"fmt"
	"log"
	d "prplchat/src/db"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
)

func HandleChatAdd(app *state.State, db *d.DBConn, user *a.User, chatName string) (string, error) {
	chat, err := db.AddChat(&d.Chat{Title: chatName, OwnerId: user.Id})
	if err != nil {
		return "", fmt.Errorf("failed to add chat to db: %s", err)
	}
	err = app.AddChat(chat.Id, chat.Title, user)
	if err != nil {
		return "", fmt.Errorf("failed to add chat[%d] to app", chat.Id)
	}
	openChat, err := app.OpenChat(user.Id, chat.Id)
	if err != nil {
		return "", fmt.Errorf("failed to open new chat: %s", err)
	}
	err = sse.DistributeChat(app, openChat, user, user, user, event.ChatAdd)
	if err != nil {
		log.Printf("HandleChatAdd ERROR cannot distribute chat[%d] creation to user[%d]: %s",
			openChat.Id, user.Id, err.Error())
	}
	template := openChat.Template(user, user)
	return template.HTML()
}

func HandleChatOpen(app *state.State, db *d.DBConn, user *a.User, chatId uint) (string, error) {
	var html string
	openChat, err := app.OpenChat(user.Id, uint(chatId))
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		welcome := template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
		html, err = welcome.HTML()
	} else {
		html, err = openChat.Template(user, user).HTML()
	}
	if err != nil {
		return "", fmt.Errorf("failed to template chat")
	}
	return html, nil
}

func HandleChatClose(app *state.State, db *d.DBConn, user *a.User, chatId uint) (string, error) {
	err := app.CloseChat(user.Id, chatId)
	if err != nil {
		return "", fmt.Errorf("close chat[%d] for user[%d]: %s", chatId, user.Id, err)
	}
	welcome := template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
	return welcome.HTML()
}

func HandleChatDelete(app *state.State, db *d.DBConn, user *a.User, chatId uint) error {
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil || chat == nil {
		return fmt.Errorf("chat[%d] inaccessible for user[%d], %s", chatId, user.Id, err)
	}
	err = app.DeleteChat(user.Id, chat)
	if err != nil {
		return fmt.Errorf("error deleting chat in app: %s", err.Error())
	}
	err = db.DeleteChat(chatId)
	if err != nil {
		return fmt.Errorf("error deleting chat in db: %s", err.Error())
	}
	err = sse.DistributeChat(app, chat, user, nil, user, event.ChatClose)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat close, %s", err.Error())
	}
	err = sse.DistributeChat(app, chat, user, nil, user, event.ChatDrop)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat deleted, %s", err.Error())
	}
	err = sse.DistributeChat(app, chat, user, nil, nil, event.ChatExpel)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat user expel, %s", err.Error())
	}
	return nil
}