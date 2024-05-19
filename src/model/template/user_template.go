package template

import (
	"bytes"
	"html/template"

	"go.chat/src/model/event"
)

type UserTemplate struct {
	ChatId      int
	ChatOwnerId uint
	UserId      uint
	UserName    string
	UserEmail   string
	//UserStatus  string
	ViewerId uint
}

func (c *UserTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(-27, c.UserId, -28)
}

func (c *UserTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.UserId, -29)
}

func (c *UserTemplate) ChatLeaveEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.UserId, -30)
}

func (c *UserTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles(
		"static/html/user_div.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
