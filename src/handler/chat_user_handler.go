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
)

func HandleUserInvite(
	state *state.State,
	db *d.DBConn,
	user *a.User,
	chatId uint,
	inviteeName string,
) (*a.Chat, *a.User, error) {
	appInvitee, err := FindUser(db, inviteeName)
	if err != nil {
		log.Printf("HandleUserInvite ERROR invitee not found [%s], %s\n", inviteeName, err.Error())
		return nil, nil, fmt.Errorf("invitee not found")
	} else if appInvitee == nil {
		log.Printf("HandleUserInvite WARN invitee not found [%s]\n", inviteeName)
		return nil, nil, nil
	}

	appChat, err := GetChat(state, db, user, chatId)
	if err != nil {
		log.Printf("HandleUserInvite ERROR user[%d] cannot invite into chat[%d], %s\n",
			user.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("cannot find chat: %s", err.Error())
	} else if appChat == nil {
		log.Printf("HandleUserInvite WARN user[%d] cannot invite into chat[%d]\n", user.Id, chatId)
		return nil, nil, fmt.Errorf("chat not found")
	}

	err = db.AddChatUser(chatId, appInvitee.Id)
	if err != nil {
		log.Printf("HandleUserInvite ERROR failed to add user[%d] to chat[%d] in db, %s\n",
			appInvitee.Id, chatId, err.Error())
		return nil, nil, fmt.Errorf("failed to add user to chat in db")
	}

	err = sse.DistributeChat(state, db, appChat, user, appInvitee, appInvitee, event.ChatInvite)
	if err != nil {
		log.Printf("HandleUserInvite WARN cannot distribute chat invite, %s\n", err.Error())
	}

	return appChat, appInvitee, nil
}

func HandleUserExpelled(state *state.State, db *d.DBConn, user *a.User, chatId uint, expelledId uint) (*a.User, error) {
	appExpelled, err := ExpelUser(state, db, user, chatId, uint(expelledId))
	if err != nil {
		log.Printf("HandleUserExpelled ERROR failed to expell, %s\n", err.Error())
		return nil, fmt.Errorf("failed to expell user, %s", err.Error())
	}
	chat, err := GetChat(state, db, user, chatId)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return nil, fmt.Errorf("failed to expell user: %s", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, user, appExpelled, appExpelled, event.ChatClose)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat close, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat close: %s", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, user, appExpelled, appExpelled, event.ChatDrop)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat deleted, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat deleted: %s", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, user, nil, appExpelled, event.ChatExpel)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat expel, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat expel: %s", err.Error())
	}
	return appExpelled, nil
}

func HandleUserLeaveChat(state *state.State, db *d.DBConn, user *a.User, chatId uint) error {
	chat, err := GetChat(state, db, user, chatId)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return fmt.Errorf("failed to leave chat: %s", err.Error())
	}
	log.Printf("HandleUserLeaveChat TRACE removing[%d] from chat[%d]\n", user.Id, chat.Id)
	if user.Id == chat.OwnerId {
		log.Printf("HandleUserLeaveChat ERROR cannot leave chat[%d] as owner\n", chatId)
		return fmt.Errorf("creator cannot leave chat")
	}
	expelled, err := ExpelUser(state, db, user, chatId, user.Id)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR failed to expell, %s\n", err.Error())
		return fmt.Errorf("failed to self-expel from chat: %s", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, expelled, expelled, expelled, event.ChatClose)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat close, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, expelled, expelled, expelled, event.ChatDrop)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat deleted, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, db, chat, expelled, nil, expelled, event.ChatLeave)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat user drop, %s\n", err.Error())
	}
	return nil
}

func GetChatUsers(db *d.DBConn, chatId uint) ([]*a.User, error) {
	dbChatUsers, err := db.GetChatUsers(chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get users of chat[%d], %s", chatId, err.Error())
	}
	appChatUsers := make([]*a.User, 0)
	for _, dbUser := range dbChatUsers {
		appChatUsers = append(appChatUsers, convert.UserDBToApp(&dbUser))
	}
	return appChatUsers, nil
}

func ExpelUser(state *state.State, db *d.DBConn, user *a.User, chatId uint, expelledId uint) (*a.User, error) {
	log.Printf("ExpelUser TRACE expelling[%d] from chat[%d]\n", expelledId, chatId)
	// veryfy user can only either leave themselves or be expelled by the owner
	if user.Id != expelledId {
		chat, err := GetChat(state, db, user, chatId)
		if err != nil {
			log.Printf("ExpelUser ERROR user[%d] cannot find chat[%d], %s\n", user.Id, chatId, err.Error())
			return nil, fmt.Errorf("user cannot find chat, %s", err.Error())
		}
		if user.Id != chat.OwnerId {
			log.Printf("ExpelUser ERROR user[%d] cannot expel user[%d] from chat[%d]\n", user.Id, expelledId, chatId)
			return nil, fmt.Errorf("failed to expel user from chat")
		}
	}
	dbExpelled, err := db.GetUser(expelledId)
	if err != nil || dbExpelled == nil {
		return nil, fmt.Errorf("user[%d] not found in db", expelledId)
	}
	err = db.RemoveChatUser(expelledId, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to remove user[%d] from chat[%d]: %s", expelledId, chatId, err.Error())
	}
	err = state.CloseChat(expelledId, chatId)
	if err != nil {
		log.Printf("ExpelUser TRACE user[%d] did not have chat[%d] open: %s", expelledId, chatId, err.Error())
	}
	appExpelled := convert.UserDBToApp(dbExpelled)
	return appExpelled, nil
}
