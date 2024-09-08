package convert

import (
	"neon-chat/src/db"
	"neon-chat/src/model/app"
)

func UserAppToDB(user *app.User) *db.User {
	return &db.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	}
}

func UserDBToApp(user *db.User, avatar *db.Avatar) *app.User {
	u := &app.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   app.UserType(user.Type),
		Status: app.UserStatus(user.Status),
		Salt:   user.Salt,
		Avatar: nil,
	}
	if avatar != nil {
		u.Avatar = AvatarDBToApp(avatar)
	}
	return u
}
