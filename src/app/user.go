package app

import (
	"neon-chat/src/app/enum"
	"neon-chat/src/template"
)

type User struct {
	Id     uint
	Name   string
	Email  string
	Type   enum.UserType
	Status enum.UserStatus
	Salt   string
	Avatar *Avatar
}

func (user *User) Template(
	chatId uint,
	chatOwnerId uint,
	viewerId uint,
) template.UserTemplate {
	return template.UserTemplate{
		ChatId:      chatId,
		ChatOwnerId: chatOwnerId,
		UserId:      user.Id,
		UserName:    user.Name,
		UserEmail:   user.Email,
		//UserStatus:  string(user.Status),
		ViewerId: viewerId,
	}
}
