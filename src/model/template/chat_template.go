package template

import (
	"bytes"
	"html/template"

	"go.chat/src/model/event"
)

type ChatTemplate struct {
	ChatId   int
	ChatName string
	User     UserTemplate
	Viewer   UserTemplate
	Owner    UserTemplate
	Users    []UserTemplate
	Messages []MessageTemplate
}

func (c *ChatTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(-98, c.User.UserId, -97)
}

func (c *ChatTemplate) ChatDropEvent() string {
	return event.ChatDrop.FormatEventName(c.ChatId, c.User.UserId, -96)
}

func (c *ChatTemplate) ChatCloseEvent() string {
	return event.ChatClose.FormatEventName(c.ChatId, c.User.UserId, -95)
}

func (c *ChatTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.User.UserId, -94)
}

func (c *ChatTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(c.ChatId, c.User.UserId, -93)
}

func (c *ChatTemplate) MessageAddEvent() string {
	return event.MessageAdd.FormatEventName(c.ChatId, c.User.UserId, -93)
}

func (c *ChatTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	chatTmpl := template.Must(template.ParseFiles(
		"static/html/chat/chat_div.html",
		"static/html/user_div.html",
		"static/html/chat/message_li.html"))
	err := chatTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *ChatTemplate) ShortHTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/nav/chat_li.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
