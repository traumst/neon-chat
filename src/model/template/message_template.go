package template

import (
	"bytes"
	"fmt"
	"html/template"
	"neon-chat/src/utils"
	"sync"
)

type MessageTemplate struct {
	mu               sync.Mutex
	IntermediateId   string
	ChatId           uint
	MsgId            uint
	Quote            *QuoteTemplate
	ViewerId         uint
	OwnerId          uint
	AuthorId         uint
	AuthorName       string
	AuthorAvatar     AvatarTemplate
	Text             string
	TextIntro        string
	MessageDropEvent string
}

func (m *MessageTemplate) getIntermediateId() string {
	if m.IntermediateId == "" {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.IntermediateId = utils.RandStringBytes(5)
	}
	return m.IntermediateId
}

func (m *MessageTemplate) GetId() uint {
	return m.MsgId
}

func (m *MessageTemplate) Shorten() uint {
	return m.ChatId
}

func (m *MessageTemplate) HTML() (string, error) {
	if err := m.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles(
		"static/html/chat/message_li.html",
		"static/html/chat/message_quote_div.html",
		"static/html/avatar_div.html"))
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
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
