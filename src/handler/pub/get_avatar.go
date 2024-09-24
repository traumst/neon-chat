package pub

import (
	"fmt"
	"neon-chat/src/app"
	"neon-chat/src/convert"
	"neon-chat/src/db"

	"github.com/jmoiron/sqlx"
)

func GetAvatar(dbConn sqlx.Ext, userId uint) (*app.Avatar, error) {
	avatar, err := db.GetAvatar(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar, %s", err)
	}
	if avatar == nil {
		return nil, nil
	}
	return convert.AvatarDBToApp(avatar), nil
}
