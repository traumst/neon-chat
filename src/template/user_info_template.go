package template

import (
	"bytes"
	"fmt"
	"html/template"
)

type UserInfoTemplate struct {
	ViewerId     uint
	UserId       uint
	UserName     string
	UserEmail    string
	UserAvatar   AvatarTemplate
	RegisterDate string
}

func (uit UserInfoTemplate) HTML() (string, error) {
	if err := uit.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(
		"static/html/contacts/user_info_div.html",
		"static/html/avatar_div.html",
	))
	if err := tmpl.Execute(&buf, uit); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (ust UserInfoTemplate) validate() error {
	if ust.UserId < 1 {
		return fmt.Errorf("UserSettingsTemplate user id cannot be 0")
	}
	if ust.UserName == "" {
		return fmt.Errorf("UserSettingsTemplate user name cannot be empty")
	}
	if ust.ViewerId < 1 {
		return fmt.Errorf("UserSettingsTemplate viewer id cannot be 0")
	}
	return nil
}
