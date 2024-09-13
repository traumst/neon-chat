package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	a "neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func UpdateAvatar(
	dbConn sqlx.Ext,
	userId uint,
	filename string,
	fileBytes []byte,
	mime string,
) (*a.Avatar, error) {
	oldAvatars, err := db.GetUserAvatars(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("fail to get avatar for user[%d]", userId)
	}
	saved, err := db.AddAvatar(dbConn, userId, filename, fileBytes, mime)
	if err != nil {
		return nil, fmt.Errorf("failed to save avatar[%s], %s", filename, err.Error())
	}
	if len(oldAvatars) > 0 {
		for _, old := range oldAvatars {
			if old == nil {
				continue
			}
			err := db.DropAvatar(dbConn, old.Id)
			if err != nil {
				log.Printf("updateAvatar ERROR failed to drop old avatar[%v]", old)
			}
		}
	}
	return convert.AvatarDBToApp(saved), nil
}
