package app

import "go.chat/model/template"

type Message struct {
	ID     int
	ChatID int
	Owner  *User
	Author *User
	Text   string
}

func (m *Message) Template(viewer *User) *template.MessageTemplate {
	return &template.MessageTemplate{
		MsgID:      m.ID,
		ChatID:     m.ChatID,
		Owner:      m.Owner.Name,
		Author:     m.Author.Name,
		Text:       m.Text,
		ActiveUser: viewer.Name,
	}
}
