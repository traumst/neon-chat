package handler

import (
	"fmt"
	"log"
	d "prplchat/src/db"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
)

func GetUser(app *state.State, db *d.DBConn, userId uint) (*a.User, error) {
	log.Printf("GetUser TRACE user[%d]\n", userId)
	dbUser, err := db.GetUser(userId)
	if err != nil {
		return nil, fmt.Errorf("user[%d] not found: %s", userId, err.Error())
	}

	appUser := UserFromDB(*dbUser)
	err = app.UpdateUser(appUser.Id, appUser)
	if err != nil {
		log.Printf("GetUser ERROR failed to cache user[%d]: %s", appUser.Id, err.Error())
		return &appUser, err
	}

	return &appUser, nil
}

func FindUser(app *state.State, db *d.DBConn, userName string) (*a.User, error) {
	log.Printf("FindUser TRACE IN user[%s]\n", userName)
	dbUser, err := db.SearchUser(userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}

	appUser := UserFromDB(*dbUser)
	err = app.UpdateUser(appUser.Id, appUser)
	if err != nil {
		log.Printf("FindUser ERROR failed to cache user[%d]: %s", appUser.Id, err.Error())
		return &appUser, err
	}
	log.Printf("FindUser TRACE OUT user[%s]\n", userName)
	return &appUser, nil
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
		appUser := UserFromDB(*dbUser)
		appUsers = append(appUsers, &appUser)
	}

	log.Printf("FindUsers TRACE OUT user[%s]\n", userName)
	return appUsers, nil
}

func InviteUser(app *state.State, db *d.DBConn, userId uint, chatId uint, inviteeId uint) error {
	log.Printf("InviteUser TRACE inviting[%d] to chat[%d]\n", inviteeId, chatId)
	appInvitee, err := app.GetUser(inviteeId)
	if err != nil || appInvitee == nil {
		log.Printf("InviteUser DEBUG refreshing user[%d] from db: %s", inviteeId, err.Error())
		dbUser, _ := db.GetUser(inviteeId)
		if dbUser == nil {
			return fmt.Errorf("user[%d] not found in db", inviteeId)
		}
	}
	err = app.InviteUser(userId, chatId, appInvitee)
	if err != nil {
		return fmt.Errorf("inviting user[%d] to chat[%d] in app: %s", appInvitee.Id, chatId, err.Error())
	}
	err = db.AddChatUser(inviteeId, chatId)
	if err != nil {
		return fmt.Errorf("failed to add user[%d] to chat[%d]: %s", inviteeId, chatId, err.Error())
	}
	return nil
}

func ExpelUser(app *state.State, db *d.DBConn, user *a.User, chatId uint, expelledId uint) (*a.User, error) {
	log.Printf("ExpelUser TRACE expelling[%d] from chat[%d]\n", expelledId, chatId)
	dbExpelled, err := db.GetUser(uint(expelledId))
	if err != nil || dbExpelled == nil {
		return nil, fmt.Errorf("user[%d] not found in db", expelledId)
	}
	err = db.RemoveChatUser(dbExpelled.Id, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to remove user[%d] from chat[%d]: %s", dbExpelled.Id, chatId, err.Error())
	}
	err = app.ExpelFromChat(dbExpelled.Id, chatId, dbExpelled.Id)
	if err != nil {
		return nil, fmt.Errorf("removing user[%d] from chat[%d]: %s", dbExpelled.Id, chatId, err.Error())
	}
	appExpelled := UserFromDB(*dbExpelled)
	return &appExpelled, nil
}

func UpdateUser(app *state.State, db *d.DBConn, template *a.User) (*a.User, error) {
	current, err := GetUser(app, db, template.Id)
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

	return current, app.UpdateUser(current.Id, *current)
}
