package template

import (
	"bytes"
	"html/template"

	"prplchat/src/model/event"
)

type UserTemplate struct {
	ChatId      uint
	ChatOwnerId uint
	UserId      uint
	UserName    string
	UserEmail   string
	//UserStatus  string
	ViewerId uint
}

func (c *UserTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(0, c.UserId, 0)
}

func (c *UserTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.UserId, 0)
}

func (c *UserTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(c.ChatId, c.UserId, 0)
}

func (c *UserTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/user_div.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *UserTemplate) ShortHTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/search/user_option.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
