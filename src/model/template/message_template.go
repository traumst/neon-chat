package template

import (
	"bytes"
	"fmt"
	"html/template"
)

type MessageTemplate struct {
	MsgId  uint
	ChatId uint
	// TODO use ids instead of names
	Owner            string
	Author           string
	Text             string
	ActiveUser       string
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
	if m.MsgId < 1 {
		return fmt.Errorf("MessageTemplate requires MsgId, but is [%d]", m.MsgId)
	}
	if m.ChatId < 1 {
		return fmt.Errorf("MessageTemplate requires ChatId, but is [%d]", m.ChatId)
	}
	// if len(m.Owner) < 1 {
	// 	return fmt.Errorf("MessageTemplate requires Owner, but is [%s]", m.Owner)
	// }
	if len(m.Author) < 1 {
		return fmt.Errorf("MessageTemplate requires Author, but is [%s]", m.Author)
	}
	if len(m.Text) < 1 {
		return fmt.Errorf("MessageTemplate requires Text, but is [%s]", m.Text)
	}
	if len(m.ActiveUser) < 1 {
		return fmt.Errorf("MessageTemplate requires ActiveUser, but is [%s]", m.ActiveUser)
	}
	if len(m.MessageDropEvent) < 1 {
		return fmt.Errorf("MessageTemplate requires MessageDropEvent, but is [%s]", m.MessageDropEvent)
	}
	return nil
}
