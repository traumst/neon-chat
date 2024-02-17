package model

import (
	"bytes"
	"html/template"
)

type MessageTemplate struct {
	ID         int
	Author     string
	Text       string
	ActiveUser string
}

func (m *MessageTemplate) GetHTML() (string, error) {
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles("html/message.html"))
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
