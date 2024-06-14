package app

import (
	"prplchat/src/model/event"
	"prplchat/src/model/template"
)

type Message struct {
	Id     int
	ChatId int
	Owner  *User
	Author *User
	Text   string
}

func (m *Message) Template(viewer *User) *template.MessageTemplate {
	return &template.MessageTemplate{
		MsgId:            m.Id,
		ChatId:           m.ChatId,
		Owner:            m.Owner.Name,
		Author:           m.Author.Name,
		Text:             m.Text,
		ActiveUser:       viewer.Name,
		MessageDropEvent: event.MessageDrop.FormatEventName(m.ChatId, m.Author.Id, m.Id),
	}
}
