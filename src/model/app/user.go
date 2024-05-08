package app

import (
	"go.chat/src/model/template"
)

type UserType string

// TODO add flags/permissions mapping

const (
	UserTypeFree UserType = "user-type-free"
)

type User struct {
	Id   uint
	Name string
	Type UserType
	Salt string
}

func (user *User) Template(
	chatId int,
	chatOwnerId uint,
	viewerId uint,
) *template.UserTemplate {
	return &template.UserTemplate{
		ChatId:      chatId,
		ChatOwnerId: chatOwnerId,
		UserId:      user.Id,
		UserName:    user.Name,
		ViewerId:    viewerId,
	}
}
