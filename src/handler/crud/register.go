package crud

import (
	"fmt"
	"log"

	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func Register(
	db *d.DBConn,
	newUser *a.User,
	pass string,
	authType a.AuthType,
) (*a.User, *a.Auth, error) {
	log.Printf("Register TRACE IN user\n")
	if db == nil || newUser == nil {
		return nil, nil, fmt.Errorf("missing mandatory args user[%v] db[%v]", newUser, db)
	}
	if newUser.Name == "" || pass == "" || newUser.Salt == "" {
		return nil, nil, fmt.Errorf("invalid args user[%s] salt[%s]", newUser.Name, newUser.Salt)
	}
	var appUser *a.User
	var err error
	if newUser.Id != 0 {
		appUser = newUser
	} else {
		appUser, err = CreateUser(db.Tx, newUser)
		if err != nil || appUser == nil {
			return nil, nil, fmt.Errorf("failed to create user[%v], %s", newUser, err)
		} else {
			log.Printf("Register TRACE user[%s] created\n", appUser.Name)
		}
	}
	auth, err := CreateAuth(db.Tx, appUser, pass, authType)
	if err != nil || auth == nil {
		return nil, nil, fmt.Errorf("failed to create auth[%s] for user[%v], %s", authType, appUser, err)
	}
	log.Printf("Register TRACE user[%d] auth[%v] created\n", appUser.Id, auth)
	return appUser, auth, nil
}
