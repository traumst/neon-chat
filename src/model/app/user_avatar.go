package app

import (
	"encoding/base64"

	t "go.chat/src/model/template"
)

type UserAvatar struct {
	Id     int
	UserId uint
	Title  string
	Size   string
	Image  []byte
	Mime   string
}

func (avatar *UserAvatar) Base64() string {
	return base64.StdEncoding.EncodeToString(avatar.Image)
}

func (a *UserAvatar) Template(viewer *User) *t.AvatarTemplate {
	return &t.AvatarTemplate{
		Id:     a.Id,
		Title:  a.Title,
		UserId: a.UserId,
		Size:   a.Size,
		Image:  a.Image,
		Mime:   a.Mime,
	}
}
