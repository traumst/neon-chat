package app

import "go.chat/model/template"

type Message struct {
	ID     int
	ChatID int
	Owner  string
	Author string
	Text   string
}

func (m *Message) Template(user string) *template.MessageTemplate {
	return &template.MessageTemplate{
		MsgID:      m.ID,
		ChatID:     m.ChatID,
		Owner:      m.Owner,
		Author:     m.Author,
		Text:       m.Text,
		ActiveUser: user,
	}
}
