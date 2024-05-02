package template

import (
	"bytes"
	"html/template"
)

type UserSettingsTemplate struct {
	UserId     uint
	ActiveUser string
}

func (h *UserSettingsTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles(
		"static/html/user_settings_div.html"))
	err := tmpl.Execute(&buf, h)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
