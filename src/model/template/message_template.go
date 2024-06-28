package template

import (
	"bytes"
	"html/template"
)

type MessageTemplate struct {
	MsgId            uint
	ChatId           uint
	Owner            string
	Author           string
	Text             string
	ActiveUser       string
	MessageDropEvent string
}

func (m *MessageTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles(
		"static/html/chat/message_li.html"))
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
