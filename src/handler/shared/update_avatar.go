package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func UpdateAvatar(
	db *d.DBConn,
	userId uint,
	filename string,
	fileBytes []byte,
	mime string,
) (*a.Avatar, error) {
	oldAvatars, err := db.GetUserAvatars(userId)
	if err != nil {
		return nil, fmt.Errorf("fail to get avatar for user[%d]", userId)
	}
	saved, err := db.AddAvatar(userId, filename, fileBytes, mime)
	if err != nil {
		return nil, fmt.Errorf("failed to save avatar[%s], %s", filename, err.Error())
	}
	if len(oldAvatars) > 0 {
		for _, old := range oldAvatars {
			if old == nil {
				continue
			}
			err := db.DropAvatar(old.Id)
			if err != nil {
				log.Printf("updateAvatar ERROR failed to drop old avatar[%v]", old)
			}
		}
	}
	return convert.AvatarDBToApp(saved), nil
}
