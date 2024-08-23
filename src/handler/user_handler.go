package handler

import (
	"fmt"
	"log"
	"prplchat/src/convert"
	d "prplchat/src/db"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
)

func GetUser(db *d.DBConn, userId uint) (*a.User, error) {
	log.Printf("GetUser TRACE user[%d]\n", userId)
	dbUser, err := db.GetUser(userId)
	if err != nil {
		return nil, fmt.Errorf("user[%d] not found: %s", userId, err.Error())
	}
	return convert.UserDBToApp(dbUser), nil
}

func FindUser(db *d.DBConn, userName string) (*a.User, error) {
	log.Printf("FindUser TRACE IN user[%s]\n", userName)
	dbUser, err := db.SearchUser(userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}

	log.Printf("FindUser TRACE OUT user[%s]\n", userName)
	return convert.UserDBToApp(dbUser), nil
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
		appUser := convert.UserDBToApp(dbUser)
		appUsers = append(appUsers, appUser)
	}

	log.Printf("FindUsers TRACE OUT user[%s]\n", userName)
	return appUsers, nil
}

func UpdateUser(state *state.State, db *d.DBConn, template *a.User) (*a.User, error) {
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
