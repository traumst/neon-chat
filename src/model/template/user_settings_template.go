package template

import (
	"bytes"
	"fmt"
	"html/template"

	"prplchat/src/model/event"
)

type UserSettingsTemplate struct {
	ChatId      uint
	ChatOwnerId uint
	UserId      uint
	UserName    string
	ViewerId    uint
	Avatar      *AvatarTemplate
}

func (c *UserSettingsTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(0, c.UserId, 0)
}

func (c *UserSettingsTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.UserId, 0)
}

func (c *UserSettingsTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(c.ChatId, c.UserId, 0)
}

func (h *UserSettingsTemplate) HTML() (string, error) {
	if err := h.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(
		"static/html/utils/user_settings_div.html",
		"static/html/user_div.html",
		"static/html/avatar_div.html"))
	if err := tmpl.Execute(&buf, h); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (c *UserSettingsTemplate) validate() error {
	if c.ChatId < 1 {
		return fmt.Errorf("UserSettingsTemplate chat id cannot be 0")
	}
	if c.ChatOwnerId < 1 {
		return fmt.Errorf("UserSettingsTemplate chat owner id cannot be 0")
	}
	if c.UserId < 1 {
		return fmt.Errorf("UserSettingsTemplate user id cannot be 0")
	}
	if c.UserName == "" {
		return fmt.Errorf("UserSettingsTemplate user name cannot be empty")
	}
	if c.ViewerId < 1 {
		return fmt.Errorf("UserSettingsTemplate viewer id cannot be 0")
	}
	if c.Avatar == nil {
		return fmt.Errorf("UserSettingsTemplate avatar cannot be nil")
	}
	return nil
}
