package template

import (
	"bytes"
	"html/template"

	"go.chat/src/model/event"
)

type UserSettingsTemplate struct {
	ChatId      int
	ChatOwnerId uint
	UserId      uint
	UserName    string
	ViewerId    uint
	Avatar      *AvatarTemplate
}

func (c *UserSettingsTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(-27, c.UserId, -28)
}

func (c *UserSettingsTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.UserId, -29)
}

func (c *UserSettingsTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(c.ChatId, c.UserId, -30)
}

func (h *UserSettingsTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(
		"static/html/utils/user_settings_div.html",
		"static/html/user_div.html",
		"static/html/avatar_div.html"))
	err := tmpl.Execute(&buf, h)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
