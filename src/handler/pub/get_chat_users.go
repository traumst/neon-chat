package pub

import (
	"fmt"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"

	"github.com/jmoiron/sqlx"
)

func GetChatUsers(dbConn sqlx.Ext, chatId uint) ([]*app.User, error) {
	chatUserIds, err := db.GetChatUserIds(dbConn, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat user ids: %s", err)
	}
	dbAvatars, err := db.GetAvatars(dbConn, chatUserIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatars for users[%v]: %s", chatUserIds, err)
	}
	dbUsers, err := db.GetUsers(dbConn, chatUserIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get users[%v]: %s", chatUserIds, err)
	}
	var users []*app.User
	for _, dbUser := range dbUsers {
		var dbAvatar *db.Avatar
		for _, avatar := range dbAvatars {
			if avatar.UserId == dbUser.Id {
				dbAvatar = avatar
				break
			}
		}
		users = append(users, convert.UserDBToApp(&dbUser, dbAvatar))
	}

	return users, nil
}
