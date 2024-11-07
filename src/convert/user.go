package convert

import (
	"neon-chat/src/app"
	"neon-chat/src/app/enum"
	"neon-chat/src/db"
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
		Type:   enum.UserType(user.Type),
		Status: enum.UserStatus(user.Status),
		Salt:   user.Salt,
		Avatar: nil,
	}
	if avatar != nil {
		u.Avatar = AvatarDBToApp(avatar)
	} else {
		u.Avatar = &app.Avatar{
			Id:     0,
			UserId: 0,
			Title:  "",
			Size:   "",
			Image:  []byte{},
			Mime:   "",
		}
	}
	return u
}
