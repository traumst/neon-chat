package handler

import (
	"fmt"
	"log"
	d "prplchat/src/db"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
)

func HandleUserInvite(
	app *state.State,
	db *d.DBConn,
	user *a.User,
	chatId uint,
	inviteeName string,
) (*a.Chat, *a.User, error) {
	appInvitee, err := findUser(app, db, inviteeName)
	if err != nil {
		log.Printf("HandleUserInvite ERROR invitee not found [%s], %s\n", inviteeName, err.Error())
		return nil, nil, fmt.Errorf("invitee not found")
	} else if appInvitee == nil {
		log.Printf("HandleUserInvite WARN invitee not found [%s]\n", inviteeName)
		return nil, nil, nil
	}

	appChat, err := GetChat(app, db, user, chatId)
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

	err = sse.DistributeChat(app, appChat, user, appInvitee, appInvitee, event.ChatInvite)
	if err != nil {
		log.Printf("HandleUserInvite WARN cannot distribute chat invite, %s\n", err.Error())
	}

	return appChat, appInvitee, nil
}

func HandleUserExpelled(app *state.State, db *d.DBConn, user *a.User, chatId uint, expelledId uint) (*a.User, error) {
	appExpelled, err := expelUser(app, db, user, chatId, uint(expelledId))
	if err != nil {
		log.Printf("HandleUserExpelled ERROR failed to expell, %s\n", err.Error())
		return nil, fmt.Errorf("failed to expell user, %s", err.Error())
	}
	chat, err := GetChat(app, db, user, chatId)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return nil, fmt.Errorf("failed to expell user: %s", err.Error())
	}
	err = sse.DistributeChat(app, chat, user, appExpelled, appExpelled, event.ChatClose)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat close, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat close: %s", err.Error())
	}
	err = sse.DistributeChat(app, chat, user, appExpelled, appExpelled, event.ChatDrop)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat deleted, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat deleted: %s", err.Error())
	}
	err = sse.DistributeChat(app, chat, user, nil, appExpelled, event.ChatExpel)
	if err != nil {
		log.Printf("HandleUserExpelled ERROR cannot distribute chat expel, %s\n", err.Error())
		return appExpelled, fmt.Errorf("cannot distribute chat expel: %s", err.Error())
	}
	return appExpelled, nil
}

func HandleUserLeaveChat(app *state.State, db *d.DBConn, user *a.User, chatId uint) error {
	chat, err := GetChat(app, db, user, chatId)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot find chat[%d], %s\n", chatId, err.Error())
		return fmt.Errorf("failed to leave chat: %s", err.Error())
	}
	log.Printf("HandleUserLeaveChat TRACE removing[%d] from chat[%d]\n", user.Id, chat.Id)
	if user.Id == chat.OwnerId {
		log.Printf("HandleUserLeaveChat ERROR cannot leave chat[%d] as owner\n", chatId)
		return fmt.Errorf("creator cannot leave chat")
	}
	expelled, err := expelUser(app, db, user, chatId, user.Id)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR failed to expell, %s\n", err.Error())
		return fmt.Errorf("failed to self-expel from chat: %s", err.Error())
	}
	err = sse.DistributeChat(app, chat, expelled, expelled, expelled, event.ChatClose)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat close, %s\n", err.Error())
	}
	err = sse.DistributeChat(app, chat, expelled, expelled, expelled, event.ChatDrop)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat deleted, %s\n", err.Error())
	}
	err = sse.DistributeChat(app, chat, expelled, nil, expelled, event.ChatLeave)
	if err != nil {
		log.Printf("HandleUserLeaveChat ERROR cannot distribute chat user drop, %s\n", err.Error())
	}
	return nil
}

func GetUser(db *d.DBConn, userId uint) (*a.User, error) {
	log.Printf("GetUser TRACE user[%d]\n", userId)
	dbUser, err := db.GetUser(userId)
	if err != nil {
		return nil, fmt.Errorf("user[%d] not found: %s", userId, err.Error())
	}
	return UserDBToApp(dbUser), nil
}

func findUser(app *state.State, db *d.DBConn, userName string) (*a.User, error) {
	log.Printf("FindUser TRACE IN user[%s]\n", userName)
	dbUser, err := db.SearchUser(userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}

	log.Printf("FindUser TRACE OUT user[%s]\n", userName)
	return UserDBToApp(dbUser), nil
}

func FindUsers(db *d.DBConn, userName string) ([]*a.User, error) {
	log.Printf("FindUsers TRACE user[%s]\n", userName)
	dbUsers, err := db.SearchUsers(userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}

	appUsers := make([]*a.User, 0)
	for _, dbUser := range dbUsers {
		if dbUser == nil {
			continue
		}
		appUser := UserDBToApp(dbUser)
		appUsers = append(appUsers, appUser)
	}

	log.Printf("FindUsers TRACE OUT user[%s]\n", userName)
	return appUsers, nil
}

func expelUser(app *state.State, db *d.DBConn, user *a.User, chatId uint, expelledId uint) (*a.User, error) {
	log.Printf("ExpelUser TRACE expelling[%d] from chat[%d]\n", expelledId, chatId)
	if user.Id == expelledId {
		return nil, fmt.Errorf("user[%d] cannot expel itself", user.Id)
	}
	dbExpelled, err := db.GetUser(expelledId)
	if err != nil || dbExpelled == nil {
		return nil, fmt.Errorf("user[%d] not found in db", expelledId)
	}
	err = db.RemoveChatUser(expelledId, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to remove user[%d] from chat[%d]: %s", expelledId, chatId, err.Error())
	}
	err = app.CloseChat(expelledId, chatId)
	appExpelled := UserDBToApp(dbExpelled)
	return appExpelled, nil
}

func UpdateUser(app *state.State, db *d.DBConn, template *a.User) (*a.User, error) {
	current, err := GetUser(db, template.Id)
	if err != nil {
		return nil, fmt.Errorf("user for update[%d] not found: %s", template.Id, err.Error())
	}

	if current.Status != template.Status {
		current.Status = template.Status
		err = db.UpdateUserStatus(template.Id, string(template.Status))
		if err != nil {
			log.Printf("UpdateUser ERROR failed to update user[%d] status: %s", template.Id, err.Error())
		}
	}
	if current.Name != template.Name {
		current.Name = template.Name
		err = db.UpdateUserName(template.Id, template.Name)
		if err != nil {
			log.Printf("UpdateUser ERROR failed to update user[%d] name: %s", template.Id, err.Error())
		}
	}

	return current, nil
}
