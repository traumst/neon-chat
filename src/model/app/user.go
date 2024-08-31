package app

import (
	t "neon-chat/src/model/template"
)

// TODO add flags/permissions mapping
type UserType string

const (
	UserTypeBasic UserType = "basic"
)

// TODO allow user ban / suspend
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
) t.UserTemplate {
	return t.UserTemplate{
		ChatId:      chatId,
		ChatOwnerId: chatOwnerId,
		UserId:      user.Id,
		UserName:    user.Name,
		UserEmail:   user.Email,
		//UserStatus:  string(user.Status),
		ViewerId: viewerId,
	}
}
