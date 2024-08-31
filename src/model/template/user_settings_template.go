package template

import (
	"bytes"
	"fmt"
	"html/template"

	ti "neon-chat/src/interface"
	"neon-chat/src/model/event"
)

type UserSettingsTemplate struct {
	ChatId      uint
	ChatOwnerId uint
	UserId      uint
	UserName    string
	ViewerId    uint
	Avatar      ti.Renderable
}

func (ust UserSettingsTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(0, ust.UserId, 0)
}

func (ust UserSettingsTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(ust.ChatId, ust.UserId, 0)
}

func (ust UserSettingsTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(ust.ChatId, ust.UserId, 0)
}

func (ust UserSettingsTemplate) HTML() (string, error) {
	if err := ust.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(
		"static/html/utils/user_settings_div.html",
		"static/html/user_div.html",
		"static/html/avatar_div.html"))
	if err := tmpl.Execute(&buf, ust); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (ust UserSettingsTemplate) validate() error {
	if ust.UserId < 1 {
		return fmt.Errorf("UserSettingsTemplate user id cannot be 0")
	}
	if ust.UserName == "" {
		return fmt.Errorf("UserSettingsTemplate user name cannot be empty")
	}
	if ust.ViewerId < 1 {
		return fmt.Errorf("UserSettingsTemplate viewer id cannot be 0")
	}
	if ust.Avatar == nil {
		return fmt.Errorf("UserSettingsTemplate avatar cannot be nil")
	}
	return nil
}
