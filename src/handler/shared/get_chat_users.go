package shared

import (
	"fmt"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func GetChatUsers(
	db *d.DBConn,
	chatId uint,
) ([]*a.User, error) {
	chatUserIds, err := db.GetChatUserIds(chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat user ids: %s", err)
	}
	dbAvatars, err := db.GetAvatars(chatUserIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatars for users[%v]: %s", chatUserIds, err)
	}
	dbUsers, err := db.GetUsers(chatUserIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get users[%v]: %s", chatUserIds, err)
	}
	var users []*a.User
	for _, dbUser := range dbUsers {
		var dbAvatar *d.Avatar
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
