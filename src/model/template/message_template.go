package template

import (
	"bytes"
	"fmt"
	"html/template"
)

type MessageTemplate struct {
	ChatId           uint
	MsgId            uint
	ViewerId         uint
	OwnerId          uint
	AuthorId         uint
	AuthorName       string
	Text             string
	MessageDropEvent string
}

func (m MessageTemplate) HTML() (string, error) {
	if err := m.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles("static/html/chat/message_li.html"))
	if err := msgTmpl.Execute(&buf, m); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (m MessageTemplate) validate() error {
	if m.ChatId < 1 {
		return fmt.Errorf("MessageTemplate requires ChatId but is [%d]", m.ChatId)
	}
	if m.MsgId < 1 {
		return fmt.Errorf("MessageTemplate requires MsgId but is [%d]", m.MsgId)
	}
	if m.ViewerId < 1 {
		return fmt.Errorf("MessageTemplate requires ViewerId but is [%d]", m.ViewerId)
	}
	if m.OwnerId < 1 {
		return fmt.Errorf("MessageTemplate requires OwnerId but is [%d]", m.OwnerId)
	}
	if m.AuthorId < 1 {
		return fmt.Errorf("MessageTemplate requires AuthorId but is [%d]", m.AuthorId)
	}
	if len(m.AuthorName) < 1 {
		return fmt.Errorf("MessageTemplate requires AuthorName but is [%s]", m.AuthorName)
	}
	if len(m.Text) < 1 {
		return fmt.Errorf("MessageTemplate requires Text but is [%s]", m.Text)
	}
	if len(m.MessageDropEvent) < 1 {
		return fmt.Errorf("MessageTemplate requires MessageDropEvent but is [%s]", m.MessageDropEvent)
	}
	return nil
}
