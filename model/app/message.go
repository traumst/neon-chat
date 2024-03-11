package app

import "go.chat/model/template"

type Message struct {
	ID     int
	ChatID int
	Author string
	Text   string
}

func (m *Message) ToTemplate(user string) *template.MessageTemplate {
	return &template.MessageTemplate{
		MsgID:      m.ID,
		ChatID:     m.ChatID,
		Author:     m.Author,
		Text:       m.Text,
		ActiveUser: user,
	}
}
