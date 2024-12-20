package pub

import (
	"fmt"
	"log"

	"neon-chat/src/app"
	"neon-chat/src/app/enum"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/utils"
)

func AuthenticateUser(
	dbConn *db.DBConn,
	username string,
	pass string,
	authType enum.AuthType,
) (*app.User, *app.Auth, error) {
	if dbConn == nil || len(username) <= 0 || len(pass) <= 0 {
		log.Printf("ERROR Authenticate bad arguments username[%s] authType[%s]\n", username, authType)
		return nil, nil, fmt.Errorf("bad arguments")
	}
	dbUser, err := db.SearchUser(dbConn.Conn, username)
	if err != nil || dbUser == nil || dbUser.Id <= 0 || len(dbUser.Salt) <= 0 {
		log.Printf("TRACE Authenticate user[%s] not found, result[%v], %s\n", username, dbUser, err)
		return nil, nil, nil
	}
	dbAvatar, _ := db.GetAvatar(dbConn.Conn, dbUser.Id)
	if dbAvatar == nil {
		log.Printf("WARN Authenticate user[%d] has no avatar\n", dbUser.Id)
		dbAvatar = &db.Avatar{}
	}
	appUser := convert.UserDBToApp(dbUser, dbAvatar)
	if appUser.Status != enum.UserStatusActive {
		log.Printf("WARN Authenticate user[%d] status[%s] is inactive\n", dbUser.Id, dbUser.Status)
		return appUser, nil, nil
	}
	hash, err := utils.HashPassword(pass, appUser.Salt)
	if err != nil {
		log.Printf("TRACE Authenticate failed on hashing[%s] pass for user[%d], %s", hash, appUser.Id, err)
		return appUser, nil, fmt.Errorf("failed hashing pass for user[%d], %s", appUser.Id, err)
	}
	log.Printf("TRACE Authenticate user[%d] auth[%s] hash[%s]\n", appUser.Id, authType, hash)
	dbAuth, err := db.GetAuth(dbConn.Conn, string(authType), hash)
	if err != nil {
		return appUser, nil, fmt.Errorf("no auth for user[%d] hash[%s], %s", appUser.Id, hash, err)
	}
	appAuth := convert.AuthDBToApp(dbAuth)
	return appUser, appAuth, nil
}
