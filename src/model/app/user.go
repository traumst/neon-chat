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

func (m *User) Template() *template.UserTemplate {
	return &template.UserTemplate{
		Id:   m.Id,
		Name: m.Name,
		//TODO UserChangeEvent: event.UserChanged.FormatEventName(0, m.Id, 0),
	}
}
