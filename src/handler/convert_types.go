package handler

import (
	"prplchat/src/db"
	"prplchat/src/model/app"
	"prplchat/src/utils"
)

func MessageAppToDB(message *app.Message) db.Message {
	return db.Message{
		Id:       message.Id,
		ChatId:   message.ChatId,
		AuthorId: message.Author.Id,
		Text:     message.Text,
	}
}

func MessageDBToApp(message *db.Message) app.Message {
	return app.Message{
		Id:     message.Id,
		ChatId: message.ChatId,
		Author: &app.User{Id: message.AuthorId},
		Text:   message.Text,
	}
}

func ChatAppToDB(chat *app.Chat) db.Chat {
	return db.Chat{
		Id:      chat.Id,
		Title:   chat.Name,
		OwnerId: chat.Owner.Id,
	}
}

func ChatDBToApp(chat *db.Chat) app.Chat {
	return app.Chat{
		Id:    chat.Id,
		Name:  chat.Title,
		Owner: &app.User{Id: chat.OwnerId},
	}
}

func UserAppToDB(user app.User) db.User {
	return db.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	}
}

func UserDBToApp(user db.User) app.User {
	return app.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   app.UserType(user.Type),
		Status: app.UserStatus(user.Status),
		Salt:   user.Salt,
	}
}

func AuthAppToDB(auth app.Auth) db.Auth {
	return db.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   string(auth.Type),
		Hash:   auth.Hash,
	}
}

func AuthDBToApp(auth db.Auth) app.Auth {
	return app.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   app.AuthType(auth.Type),
		Hash:   auth.Hash,
	}
}

func AvatarAppToDB(avatar app.Avatar) db.Avatar {
	return db.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   int(utils.SizeDecode(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}

func AvatarDBToApp(avatar db.Avatar) app.Avatar {
	return app.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   utils.SizeEncode(int64(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}
