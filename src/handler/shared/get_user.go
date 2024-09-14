package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func GetUser(dbConn sqlx.Ext, userId uint) (*a.User, error) {
	log.Printf("GetUser TRACE user[%d]\n", userId)
	dbUser, err := db.GetUser(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("user[%d] not found: %s", userId, err.Error())
	}
	var dbAvatar *d.Avatar
	if dbUser != nil {
		dbAvatar, _ = db.GetAvatar(dbConn, dbUser.Id)
	}
	return convert.UserDBToApp(dbUser, dbAvatar), nil
}
