package handler

import (
	"prplchat/src/db"
	"prplchat/src/model/app"
	"prplchat/src/utils"
)

func UserToDB(user app.User) db.User {
	return db.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	}
}

func UserFromDB(user db.User) app.User {
	return app.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   app.UserType(user.Type),
		Status: app.UserStatus(user.Status),
		Salt:   user.Salt,
	}
}

func AuthToDB(auth app.Auth) db.Auth {
	return db.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   string(auth.Type),
		Hash:   auth.Hash,
	}
}

func AuthFromDB(auth db.Auth) app.Auth {
	return app.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   app.AuthType(auth.Type),
		Hash:   auth.Hash,
	}
}

func AvatarToDB(avatar app.Avatar) db.Avatar {
	return db.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   int(utils.SizeDecode(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}

func AvatarFromDB(avatar db.Avatar) app.Avatar {
	return app.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   utils.SizeEncode(int64(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}
