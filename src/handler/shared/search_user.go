package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func SearchUser(db *d.DBConn, userName string) (*a.User, error) {
	log.Printf("FindUser TRACE IN user[%s]\n", userName)
	dbUser, err := db.SearchUser(userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}
	var dbAvatar *d.Avatar
	if dbUser != nil {
		dbAvatar, _ = db.GetAvatar(dbUser.Id)
	}

	log.Printf("FindUser TRACE OUT user[%s]\n", userName)
	return convert.UserDBToApp(dbUser, dbAvatar), nil
}
