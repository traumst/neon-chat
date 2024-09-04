package app

import (
	"fmt"
	"neon-chat/src/model/event"
	t "neon-chat/src/model/template"
	"neon-chat/src/utils"
)

type Message struct {
	Id     uint
	ChatId uint
	Author *User
	Text   string
	Quote  *Message
}

func (m *Message) Template(
	viewer *User,
	owner *User,
	avatar *Avatar,
	quote *Message,
) (t.MessageTemplate, error) {
	if viewer == nil || viewer.Id == 0 {
		return t.MessageTemplate{}, fmt.Errorf("viewer cannot be nil or blank")
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
	var quoteTmpl *t.MessageTemplate
	if quote != nil {
		quoteTmplVal, err := quote.Template(viewer, owner, avatar, nil)
		if err != nil {
			return t.MessageTemplate{}, fmt.Errorf("failed to template quote: %s", err)
		}
		quoteTmpl = &quoteTmplVal
	}
	return t.MessageTemplate{
		IntermediateId:   utils.RandStringBytes(5),
		ChatId:           m.ChatId,
		MsgId:            m.Id,
		Quote:            quoteTmpl,
		ViewerId:         viewer.Id,
		OwnerId:          owner.Id,
		AuthorId:         m.Author.Id,
		AuthorName:       m.Author.Name,
		AuthorAvatar:     avatar.Template(viewer),
		Text:             m.Text,
		TextIntro:        utils.Shorten(m.Text, 69),
		MessageDropEvent: event.MessageDrop.FormatEventName(m.ChatId, m.Author.Id, m.Id),
	}, nil
}
