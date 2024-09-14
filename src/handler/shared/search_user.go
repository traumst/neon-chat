package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func SearchUser(dbConn sqlx.Ext, userName string) (*a.User, error) {
	log.Printf("FindUser TRACE IN user[%s]\n", userName)
	dbUser, err := d.SearchUser(dbConn, userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}
	var dbAvatar *d.Avatar
	if dbUser != nil {
		dbAvatar, _ = d.GetAvatar(dbConn, dbUser.Id)
	}

	log.Printf("FindUser TRACE OUT user[%s]\n", userName)
	return convert.UserDBToApp(dbUser, dbAvatar), nil
}
