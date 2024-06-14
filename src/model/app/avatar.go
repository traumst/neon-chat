package app

import (
	t "prplchat/src/model/template"
)

type Avatar struct {
	Id     int
	UserId uint
	Title  string
	Size   string
	Image  []byte
	Mime   string
}

func (a *Avatar) Template(viewer *User) *t.AvatarTemplate {
	return &t.AvatarTemplate{
		Id:     a.Id,
		Title:  a.Title,
		UserId: a.UserId,
		Size:   a.Size,
		Image:  a.Image,
		Mime:   a.Mime,
	}
}
