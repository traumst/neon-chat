package handler

import (
	"fmt"
	"log"

	"neon-chat/src/db"
	"neon-chat/src/model/app"
)

func RegisterUser(
	dbConn *db.DBConn,
	newUser *app.User,
	pass string,
	authType app.AuthType,
) (*app.User, *app.Auth, error) {
	log.Printf("Register TRACE IN user\n")
	if dbConn == nil || newUser == nil {
		return nil, nil, fmt.Errorf("missing mandatory args user[%v] db[%v]", newUser, dbConn)
	}
	if newUser.Name == "" || pass == "" || newUser.Salt == "" {
		return nil, nil, fmt.Errorf("invalid args user[%s] salt[%s]", newUser.Name, newUser.Salt)
	}
	var appUser *app.User
	var err error
	if newUser.Id != 0 {
		appUser = newUser
	} else {
		appUser, err = CreateUser(dbConn.Tx, newUser)
		if err != nil || appUser == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", newUser, err)
		} else {
			log.Printf("Register TRACE user[%s] created\n", appUser.Name)
		}
	}
	auth, err := CreateAuth(dbConn.Tx, appUser, pass, authType)
	if err != nil || auth == nil {
		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, appUser, err)
	}
	log.Printf("Register TRACE user[%d] auth[%v] created\n", appUser.Id, auth)
	return appUser, auth, nil
}
