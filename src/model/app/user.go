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
	Id   uint     `db:"id"`
	Name string   `db:"name"`
	Type UserType `db:"type"`
	Salt string   `db:"salt"`
}

func (m *User) Template(
	chatId int,
	chatOwnerId uint,
	viewerId uint,
) *template.UserTemplate {
	return &template.UserTemplate{
		ChatId:      chatId,
		ChatOwnerId: chatOwnerId,
		UserId:      m.Id,
		UserName:    m.Name,
		ViewerId:    viewerId,
	}
}
