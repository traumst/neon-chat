package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func GetUser(db *d.DBConn, userId uint) (*a.User, error) {
	log.Printf("GetUser TRACE user[%d]\n", userId)
	dbUser, err := db.GetUser(userId)
	if err != nil {
		return nil, fmt.Errorf("user[%d] not found: %s", userId, err.Error())
	}
	var dbAvatar *d.Avatar
	if dbUser != nil {
		dbAvatar, _ = db.GetAvatar(dbUser.Id)
	}
	return convert.UserDBToApp(dbUser, dbAvatar), nil
}
