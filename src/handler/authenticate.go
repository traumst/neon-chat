package handler

import (
	"fmt"
	"log"

	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
	"neon-chat/src/utils"
)

func Authenticate(
	db *d.DBConn,
	username string,
	pass string,
	authType a.AuthType,
) (*a.User, *a.Auth, error) {
	if db == nil || len(username) <= 0 || len(pass) <= 0 {
		log.Printf("Authenticate ERROR bad arguments username[%s] authType[%s]\n", username, authType)
		return nil, nil, fmt.Errorf("bad arguments")
	}
	dbUser, err := d.SearchUser(db.Conn, username)
	if err != nil || dbUser == nil || dbUser.Id <= 0 || len(dbUser.Salt) <= 0 {
		log.Printf("Authenticate TRACE user[%s] not found, result[%v], %s\n", username, dbUser, err)
		return nil, nil, nil
	}
	dbAvatar, _ := d.GetAvatar(db.Conn, dbUser.Id)
	appUser := convert.UserDBToApp(dbUser, dbAvatar)
	if appUser.Status != a.UserStatusActive {
		log.Printf("Authenticate WARN user[%d] status[%s] is inactive\n", dbUser.Id, dbUser.Status)
		return appUser, nil, nil
	}
	hash, err := utils.HashPassword(pass, appUser.Salt)
	if err != nil {
		log.Printf("Authenticate TRACE failed on hashing[%s] pass for user[%d], %s", hash, appUser.Id, err)
		return appUser, nil, fmt.Errorf("failed hashing pass for user[%d], %s", appUser.Id, err)
	}
	log.Printf("Authenticate TRACE user[%d] auth[%s] hash[%s]\n", appUser.Id, authType, hash)
	dbAuth, err := d.GetAuth(db.Conn, string(authType), hash)
	if err != nil {
		return appUser, nil, fmt.Errorf("no auth for user[%d] hash[%s], %s", appUser.Id, hash, err)
	}
	appAuth := convert.AuthDBToApp(dbAuth)
	return appUser, appAuth, nil
}
