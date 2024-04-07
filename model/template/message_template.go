package template

import (
	"bytes"
	"html/template"
)

type MessageTemplate struct {
	MsgId      int
	ChatId     int
	Owner      string
	Author     string
	Text       string
	ActiveUser string
}

func (m *MessageTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles("html/bits/message_li.html"))
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
