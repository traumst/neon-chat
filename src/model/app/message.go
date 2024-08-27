package app

import (
	"fmt"
	"prplchat/src/model/event"
	t "prplchat/src/model/template"
)

type Message struct {
	Id     uint
	ChatId uint
	Author *User
	Text   string
}

func (m *Message) Template(viewer *User, owner *User, avatar *Avatar) (t.MessageTemplate, error) {
	if viewer == nil || viewer.Id == 0 {
		return t.MessageTemplate{}, fmt.Errorf("viewer cannot be nil or blank")
	}
	if owner == nil || owner.Id == 0 {
		return t.MessageTemplate{}, fmt.Errorf("owner cannot be nil or blank")
	}
	if m.Author == nil || m.Author.Id == 0 || m.Author.Name == "" {
		return t.MessageTemplate{}, fmt.Errorf("message author cannot be nil or blank")
	}
	if m.ChatId == 0 {
		return t.MessageTemplate{}, fmt.Errorf("message chatId cannot be 0")
	}
	if m.Id == 0 {
		return t.MessageTemplate{}, fmt.Errorf("message chatId and Id cannot be 0")
	}
	return t.MessageTemplate{
		ChatId:           m.ChatId,
		MsgId:            m.Id,
		ViewerId:         viewer.Id,
		OwnerId:          owner.Id,
		AuthorId:         m.Author.Id,
		AuthorName:       m.Author.Name,
		AuthorAvatar:     avatar.Template(viewer),
		Text:             m.Text,
		MessageDropEvent: event.MessageDrop.FormatEventName(m.ChatId, m.Author.Id, m.Id),
	}, nil
}
