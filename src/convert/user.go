package convert

import (
	"prplchat/src/db"
	"prplchat/src/model/app"
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

func UserDBToApp(user *db.User) *app.User {
	return &app.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   app.UserType(user.Type),
		Status: app.UserStatus(user.Status),
		Salt:   user.Salt,
		Avatar: nil,
	}
}
