package template

import (
	"bytes"
	"fmt"
	"html/template"
	"neon-chat/src/utils"
	"sync"
)

type QuoteTemplate struct {
	mu             sync.Mutex
	IntermediateId string
	ChatId         uint
	MsgId          uint
	AuthorId       uint
	AuthorName     string
	AuthorAvatar   AvatarTemplate
	Text           string
	TextIntro      string
}

func (m *QuoteTemplate) getIntermediateId() string {
	if m.IntermediateId == "" {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.IntermediateId = utils.RandStringBytes(5)
	}
	return m.IntermediateId
}

func (m *QuoteTemplate) GetId() uint {
	return m.MsgId
}

func (m *QuoteTemplate) HTML() (string, error) {
	if err := m.validate(); err != nil {
		return "", fmt.Errorf("failed to validate QuoteTemplate, %s", err.Error())
	}
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles(
		"static/html/chat/message_submit_quote_div.html",
		"static/html/avatar_div.html"))
	err := msgTmpl.Execute(&buf, QuoteTemplate{
		IntermediateId: m.getIntermediateId(),
		ChatId:         m.ChatId,
		MsgId:          m.MsgId,
		AuthorId:       m.AuthorId,
		AuthorName:     m.AuthorName,
		AuthorAvatar:   m.AuthorAvatar,
		Text:           m.Text,
		TextIntro:      m.TextIntro,
	})
	if err != nil {
		return "", fmt.Errorf("failed to short template, %s", err.Error())
	}
	return buf.String(), nil
}

func (m *QuoteTemplate) validate() error {
	if m.ChatId < 1 {
		return fmt.Errorf("MessageTemplate requires ChatId but is [%d]", m.ChatId)
	}
	if m.MsgId < 1 {
		return fmt.Errorf("MessageTemplate requires MsgId but is [%d]", m.MsgId)
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
	return nil
}
