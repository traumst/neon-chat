package app

import (
	"go.chat/src/model/template"
)

// TODO add flags/permissions mapping

type UserType string

const (
	UserTypeFree UserType = "user-type-free"
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
		UserEmail:   user.Email,
		//UserStatus:  string(user.Status),
		ViewerId: viewerId,
	}
}
