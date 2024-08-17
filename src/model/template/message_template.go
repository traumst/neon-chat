package template

import (
	"bytes"
	"fmt"
	"html/template"
)

type MessageTemplate struct {
	ChatId           uint
	MsgId            uint
	ViewerName       string
	OwnerName        string
	AuthorName       string
	Text             string
	MessageDropEvent string
}

func (m *MessageTemplate) HTML() (string, error) {
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

func (m *MessageTemplate) validate() error {
	if m.ChatId < 1 {
		return fmt.Errorf("MessageTemplate requires ChatId but is [%d]", m.ChatId)
	}
	if m.MsgId < 1 {
		return fmt.Errorf("MessageTemplate requires MsgId but is [%d]", m.MsgId)
	}
	if len(m.ViewerName) < 1 {
		return fmt.Errorf("MessageTemplate requires ViewerName but is [%s]", m.ViewerName)
	}
	if len(m.OwnerName) < 1 {
		return fmt.Errorf("MessageTemplate requires OwnerName but is [%s]", m.OwnerName)
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
