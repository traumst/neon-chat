package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/sse"
	"neon-chat/src/handler/state"
	a "neon-chat/src/model/app"
	"neon-chat/src/model/event"
)

func HandleChatAdd(state *state.State, db *d.DBConn, user *a.User, chatName string) (string, error) {
	dbChat, err := d.AddChat(db.Tx, &d.Chat{Title: chatName, OwnerId: user.Id})
	if err != nil {
		return "", fmt.Errorf("failed to add chat[%s] to db: %s", chatName, err)
	}
	err = d.AddChatUser(db.Tx, dbChat.Id, user.Id)
	if err != nil {
		return "", fmt.Errorf("failed to add owner[%d] to chat[%d] in db: %s", user.Id, dbChat.Id, err)
	}
	err = state.OpenChat(user.Id, dbChat.Id)
	if err != nil {
		return "", fmt.Errorf("failed to open new chat: %s", err)
	}
	openChat, err := d.GetChat(db.Tx, dbChat.Id)
	if err != nil {
		return "", fmt.Errorf("failed to get chat[%d] from db: %s", dbChat.Id, err)
	}
	appChat := convert.ChatDBToApp(openChat, &d.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	})
	err = sse.DistributeChat(state, db.Tx, appChat, user, user, user, event.ChatAdd)
	if err != nil {
		log.Printf("HandleChatAdd ERROR cannot distribute chat[%d] creation to user[%d]: %s",
			openChat.Id, user.Id, err.Error())
	}
	appChatUsers, err := shared.GetChatUsers(db.Tx, dbChat.Id)
	if err != nil {
		log.Printf("HandleChatAdd ERROR getting chat[%d] users: %s", dbChat.Id, err)
	}
	appChatMsgs := make([]*a.Message, 0)
	return appChat.Template(user, user, appChatUsers, appChatMsgs).HTML()
}

func HandleChatOpen(state *state.State, db *d.DBConn, user *a.User, chatId uint) (string, error) {
	err := state.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		return TemplateWelcome(user)
	}
	openChatId := state.GetOpenChat(user.Id)
	if openChatId != chatId {
		panic(fmt.Errorf("chat[%d] should have been open for user[%d]", chatId, user.Id))
	}
	appChat, err := shared.GetChat(state, db.Conn, user, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		return TemplateWelcome(user)
	}
	appChatUsers, err := shared.GetChatUsers(db.Conn, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR getting chat[%d] users for user[%d], %s\n", chatId, user.Id, err.Error())
		return TemplateWelcome(user)
	}
	appChatMsgs, err := shared.GetChatMessages(db.Conn, chatId)
	if err != nil {
		log.Printf("HandleChatAdd ERROR getting chat[%d] messages: %s", appChat.Id, err)
	}
	return appChat.Template(user, user, appChatUsers, appChatMsgs).HTML()
}

func HandleChatClose(state *state.State, db *d.DBConn, user *a.User, chatId uint) (string, error) {
	err := state.CloseChat(user.Id, chatId)
	if err != nil {
		return "", fmt.Errorf("close chat[%d] for user[%d]: %s", chatId, user.Id, err)
	}
	return TemplateWelcome(user)
}

func HandleChatDelete(state *state.State, db *d.DBConn, user *a.User, chatId uint) error {
	dbChat, err := d.GetChat(db.Tx, chatId)
	if err != nil {
		log.Printf("HandleChatDelete ERROR chat[%d] not found in db: %s", chatId, err)
		return nil
	}
	if user.Id != dbChat.OwnerId {
		return fmt.Errorf("user[%d] is not owner[%d] of chat[%d]", user.Id, dbChat.OwnerId, chatId)
	}
	err = d.DeleteChat(db.Tx, chatId)
	if err != nil {
		return fmt.Errorf("error deleting chat in db: %s", err.Error())
	}
	chat := convert.ChatDBToApp(dbChat, &d.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	})
	err = sse.DistributeChat(state, db.Tx, chat, user, nil, user, event.ChatClose)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat close, %s", err.Error())
	}
	err = sse.DistributeChat(state, db.Tx, chat, user, nil, user, event.ChatDrop)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat deleted, %s", err.Error())
	}
	err = sse.DistributeChat(state, db.Tx, chat, user, nil, nil, event.ChatExpel)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat user expel, %s", err.Error())
	}
	return nil
}
