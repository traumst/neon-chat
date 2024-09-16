package handler

import (
	"fmt"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func GetAvatar(dbConn sqlx.Ext, userId uint) (*app.Avatar, error) {
	avatar, err := db.GetAvatar(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar, %s", err)
	}
	return convert.AvatarDBToApp(avatar), nil

}
