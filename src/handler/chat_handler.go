package handler

import (
	"fmt"
	"log"
	"prplchat/src/convert"
	d "prplchat/src/db"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
	"prplchat/src/utils"
)

func HandleChatAdd(state *state.State, db *d.DBConn, user *a.User, chatName string) (string, error) {
	chat, err := db.AddChat(&d.Chat{Title: chatName, OwnerId: user.Id})
	if err != nil {
		return "", fmt.Errorf("failed to add chat[%s] to db: %s", chatName, err)
	}
	err = db.AddChatUser(chat.Id, user.Id)
	if err != nil {
		return "", fmt.Errorf("failed to add owner[%d] to chat[%d] in db: %s", user.Id, chat.Id, err)
	}
	err = state.OpenChat(user.Id, chat.Id)
	if err != nil {
		return "", fmt.Errorf("failed to open new chat: %s", err)
	}
	openChat, err := db.GetChat(chat.Id)
	if err != nil {
		return "", fmt.Errorf("failed to get chat[%d] from db: %s", chat.Id, err)
	}
	appChat := convert.ChatDBToApp(openChat)
	err = sse.DistributeChat(state, db, appChat, user, user, user, event.ChatAdd)
	if err != nil {
		log.Printf("HandleChatAdd ERROR cannot distribute chat[%d] creation to user[%d]: %s",
			openChat.Id, user.Id, err.Error())
	}
	dbChatUsers, err := db.GetChatUsers(chat.Id)
	if err != nil {
		log.Printf("HandleChatAdd ERROR getting chat[%d] users: %s", chat.Id, err)
	}
	appChatUsers := make([]*a.User, 0)
	for _, dbChatUser := range dbChatUsers {
		appChatUsers = append(appChatUsers, convert.UserDBToApp(&dbChatUser))
	}
	appChatMsgs := make([]*a.Message, 0)
	return appChat.Template(user, user, appChatUsers, appChatMsgs).HTML()
}

func HandleChatOpen(state *state.State, db *d.DBConn, user *a.User, chatId uint) (string, error) {
	err := state.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		welcome := template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
		return welcome.HTML()
	}
	appChat, err := GetChat(state, db, user, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR opening chat[%d] for user[%d], %s\n", chatId, user.Id, err.Error())
		welcome := template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
		return welcome.HTML()
	}
	appChatUsers, err := GetChatUsers(db, chatId)
	if err != nil {
		log.Printf("HandleChatOpen ERROR getting chat[%d] users for user[%d], %s\n", chatId, user.Id, err.Error())
		return appChat.Template(user, user, nil, nil).HTML()
	}
	appChatMsgs, err := GetChatMessages(db, chatId)
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
	welcome := template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
	return welcome.HTML()
}

func HandleChatDelete(state *state.State, db *d.DBConn, user *a.User, chatId uint) error {
	dbChat, err := db.GetChat(chatId)
	if err != nil {
		log.Printf("HandleChatDelete ERROR chat[%d] not found in db: %s", chatId, err)
		return nil
	}
	if user.Id != dbChat.OwnerId {
		return fmt.Errorf("user[%d] is not owner[%d] of chat[%d]", user.Id, dbChat.OwnerId, chatId)
	}
	err = db.DeleteChat(chatId)
	if err != nil {
		return fmt.Errorf("error deleting chat in db: %s", err.Error())
	}
	chat := convert.ChatDBToApp(dbChat)
	err = sse.DistributeChat(state, db, chat, user, nil, user, event.ChatClose)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat close, %s", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, user, nil, user, event.ChatDrop)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat deleted, %s", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, user, nil, nil, event.ChatExpel)
	if err != nil {
		log.Printf("HandleChatDelete ERROR cannot distribute chat user expel, %s", err.Error())
	}
	return nil
}

func GetChat(state *state.State, db *d.DBConn, user *a.User, chatId uint) (*a.Chat, error) {
	dbChat, err := db.GetChat(chatId)
	if err != nil {
		return nil, fmt.Errorf("chat[%d] not found in db: %s", chatId, err)
	}
	if dbChat == nil {
		return nil, nil
	}
	chatIds, err := db.GetUserChatIds(user.Id)
	if err != nil {
		return nil, fmt.Errorf("error getting chat ids for user[%d]: %s", user.Id, err)
	}
	if !utils.Contains(chatIds, chatId) {
		return nil, fmt.Errorf("user[%d] is not in chat[%d]", user.Id, chatId)
	}
	return convert.ChatDBToApp(dbChat), nil
}

func GetChats(state *state.State, db *d.DBConn, userId uint) ([]*a.Chat, error) {
	dbUserChats, err := db.GetUserChats(userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user chats: %s", err)
	}
	userChats := make([]*a.Chat, 0)
	for _, dbChat := range dbUserChats {
		appChatOwner, _ := GetUser(db, dbChat.OwnerId)
		if appChatOwner == nil {
			log.Printf("GetChats WARN chat[%d] owner[%d] not found\n", dbChat.Id, dbChat.OwnerId)
			continue
		}
		chat := convert.ChatDBToApp(&dbChat)
		if chat == nil {
			log.Printf("GetChats WARN chat[%d] failed to map from db to app\n", dbChat.Id)
			continue
		}
		userChats = append(userChats, chat)
	}
	return userChats, nil
}
