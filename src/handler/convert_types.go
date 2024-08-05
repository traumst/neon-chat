package handler

import (
	"prplchat/src/db"
	a "prplchat/src/model/app"
	"prplchat/src/utils"
)

func MessageAppToDB(message *a.Message) db.Message {
	return db.Message{
		Id:       message.Id,
		ChatId:   message.ChatId,
		AuthorId: message.Author.Id,
		Text:     message.Text,
	}
}

func MessageDBToApp(message *db.Message, author *a.User) a.Message {
	return a.Message{
		Id:     message.Id,
		ChatId: message.ChatId,
		Author: author,
		Text:   message.Text,
	}
}

func ChatAppToDB(chat *a.Chat) *db.Chat {
	if chat == nil {
		return nil
	}
	return &db.Chat{
		Id:      chat.Id,
		Title:   chat.Name,
		OwnerId: chat.Owner.Id,
	}
}

func ChatDBToApp(chat *db.Chat, owner *a.User) *a.Chat {
	if chat == nil {
		return nil
	}
	if owner == nil {
		return nil
	}
	return &a.Chat{
		Id:    chat.Id,
		Name:  chat.Title,
		Owner: owner,
	}
}

func UserAppToDB(user *a.User) *db.User {
	return &db.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   string(user.Type),
		Status: string(user.Status),
		Salt:   user.Salt,
	}
}

func UserDBToApp(user *db.User) *a.User {
	return &a.User{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Type:   a.UserType(user.Type),
		Status: a.UserStatus(user.Status),
		Salt:   user.Salt,
	}
}

func AuthAppToDB(auth *a.Auth) *db.Auth {
	return &db.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   string(auth.Type),
		Hash:   auth.Hash,
	}
}

func AuthDBToApp(auth *db.Auth) *a.Auth {
	return &a.Auth{
		Id:     auth.Id,
		UserId: auth.UserId,
		Type:   a.AuthType(auth.Type),
		Hash:   auth.Hash,
	}
}

func AvatarAppToDB(avatar *a.Avatar) *db.Avatar {
	return &db.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   int(utils.SizeDecode(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}

func AvatarDBToApp(avatar *db.Avatar) *a.Avatar {
	return &a.Avatar{
		Id:     avatar.Id,
		UserId: avatar.UserId,
		Title:  avatar.Title,
		Size:   utils.SizeEncode(int64(avatar.Size)),
		Image:  avatar.Image,
		Mime:   avatar.Mime,
	}
}
