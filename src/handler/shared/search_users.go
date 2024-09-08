package shared

import (
	"fmt"
	"log"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	a "neon-chat/src/model/app"
)

func SearchUsers(db *d.DBConn, userName string) ([]*a.User, error) {
	log.Printf("FindUsers TRACE user[%s]\n", userName)
	dbUsers, err := db.SearchUsers(userName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", userName, err.Error())
	}
	dbUserIds := make([]uint, 0)
	for _, dbUser := range dbUsers {
		if dbUser == nil {
			continue
		}
		dbUserIds = append(dbUserIds, dbUser.Id)
	}
	dbAvatars, err := db.GetAvatars(dbUserIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatars for users[%v]: %s", dbUserIds, err.Error())
	}
	avatarByUserId := make(map[uint]*d.Avatar)
	for _, avatar := range dbAvatars {
		avatarByUserId[avatar.UserId] = avatar
	}
	appUsers := make([]*a.User, 0)
	for _, dbUser := range dbUsers {
		if dbUser == nil {
			continue
		}
		appUser := convert.UserDBToApp(dbUser, avatarByUserId[dbUser.Id])
		appUsers = append(appUsers, appUser)
	}
	log.Printf("FindUsers TRACE OUT user[%s]\n", userName)
	return appUsers, nil
}