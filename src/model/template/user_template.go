package template

import (
	"bytes"
	"html/template"
)

type UserTemplate struct {
	Id              uint
	Name            string
	UserChangeEvent string
}

func (m *UserTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles(
		"static/html/chat/user_name.html"))
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
