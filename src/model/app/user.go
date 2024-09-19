package app

import (
	"neon-chat/src/model/template"
)

type UserType string

const (
	UserTypeBasic UserType = "basic"
)

type UserStatus string

const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
	UserStatusSuspend UserStatus = "suspend"
)

type User struct {
	Id     uint
	Name   string
	Email  string
	Type   UserType
	Status UserStatus
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
