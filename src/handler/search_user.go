package handler

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func SearchUser(dbConn sqlx.Ext, userName string) (*app.User, error) {
	log.Printf("FindUser TRACE IN user[%s]\n", userName)
	dbUser, err := db.SearchUser(dbConn, userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}
	var dbAvatar *db.Avatar
	if dbUser != nil {
		dbAvatar, _ = db.GetAvatar(dbConn, dbUser.Id)
	}

	log.Printf("FindUser TRACE OUT user[%s]\n", userName)
	return convert.UserDBToApp(dbUser, dbAvatar), nil
}
