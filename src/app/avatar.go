package app

import (
	"neon-chat/src/template"
)

type Avatar struct {
	Id     uint
	UserId uint
	Title  string
	Size   string
	Image  []byte
	Mime   string
}

func (a *Avatar) Template(viewer *User) template.AvatarTemplate {
	if a == nil || a.Id == 0 || viewer == nil || viewer.Id == 0 {
		return template.AvatarTemplate{
			Id:     0,
			Title:  "",
			UserId: 0,
			Size:   "0KB",
			Image:  []byte{},
			Mime:   "",
		}
	}
	return template.AvatarTemplate{
		Id:     a.Id,
		Title:  a.Title,
		UserId: a.UserId,
		Size:   a.Size,
		Image:  a.Image,
		Mime:   a.Mime,
	}
}
