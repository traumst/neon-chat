package shared

import (
	"fmt"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func GetAvatar(db *d.DBConn, userId uint) (*a.Avatar, error) {
	avatar, err := db.GetAvatar(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar, %s", err)
	}
	return convert.AvatarDBToApp(avatar), nil

}
