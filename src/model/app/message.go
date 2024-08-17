package app

import (
	"prplchat/src/model/event"
	"prplchat/src/model/template"
)

type Message struct {
	Id     uint
	ChatId uint
	Author *User
	Text   string
}

func (m *Message) Template(viewer *User, ownerName string) *template.MessageTemplate {
	return &template.MessageTemplate{
		MsgId:  m.Id,
		ChatId: m.ChatId,
		//Owner:            ownerName,
		Author:           m.Author.Name,
		Text:             m.Text,
		ActiveUser:       viewer.Name,
		MessageDropEvent: event.MessageDrop.FormatEventName(m.ChatId, m.Author.Id, m.Id),
	}
}
