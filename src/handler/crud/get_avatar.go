package crud

import (
	"fmt"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func GetAvatar(dbConn sqlx.Ext, userId uint) (*a.Avatar, error) {
	avatar, err := d.GetAvatar(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar, %s", err)
	}
	return convert.AvatarDBToApp(avatar), nil

}
