package app

import (
	"fmt"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
)

type Message struct {
	Id     uint
	ChatId uint
	Author *User
	Text   string
}

func (m *Message) Template(viewer *User, owner *User) (*template.MessageTemplate, error) {
	if viewer == nil || viewer.Id == 0 {
		return nil, fmt.Errorf("viewer cannot be nil or blank")
	}
	if owner == nil || owner.Id == 0 {
		return nil, fmt.Errorf("owner cannot be nil or blank")
	}
	if m.Author == nil || m.Author.Id == 0 || m.Author.Name == "" {
		return nil, fmt.Errorf("message author cannot be nil or blank")
	}
	if m.ChatId == 0 {
		return nil, fmt.Errorf("message chatId cannot be 0")
	}
	if m.Id == 0 {
		return nil, fmt.Errorf("message chatId and Id cannot be 0")
	}
	return &template.MessageTemplate{
		ChatId:           m.ChatId,
		MsgId:            m.Id,
		ViewerName:       viewer.Name,
		OwnerName:        owner.Name,
		AuthorName:       m.Author.Name,
		Text:             m.Text,
		MessageDropEvent: event.MessageDrop.FormatEventName(m.ChatId, m.Author.Id, m.Id),
	}, nil
}
