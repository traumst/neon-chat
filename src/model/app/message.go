package app

import (
	"fmt"
	"log"
	"neon-chat/src/model/event"
	"neon-chat/src/model/template"
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
	quote *Message,
) (template.MessageTemplate, error) {
	if viewer == nil || viewer.Id == 0 {
		return template.MessageTemplate{}, fmt.Errorf("viewer cannot be nil or blank")
	}
	if m.Author == nil || m.Author.Id == 0 || m.Author.Name == "" {
		return template.MessageTemplate{}, fmt.Errorf("message author cannot be nil or blank")
	}
	if m.ChatId == 0 {
		return template.MessageTemplate{}, fmt.Errorf("message chatId cannot be 0")
	}
	if m.Id == 0 {
		return template.MessageTemplate{}, fmt.Errorf("message chatId and Id cannot be 0")
	}
	var quoteTmpl template.QuoteTemplate
	if quote != nil {
		tmpl, err := quote.Template(viewer, owner, nil)
		if err != nil {
			return template.MessageTemplate{}, fmt.Errorf("failed to template quote: %s", err)
		}
		quoteTmpl = template.QuoteTemplate{
			IntermediateId: tmpl.IntermediateId,
			ChatId:         tmpl.ChatId,
			MsgId:          tmpl.MsgId,
			AuthorId:       tmpl.AuthorId,
			AuthorName:     tmpl.AuthorName,
			AuthorAvatar:   tmpl.AuthorAvatar,
			Text:           tmpl.Text,
			TextIntro:      tmpl.TextIntro,
		}
	}
	var avatarTmpl template.AvatarTemplate
	if m.Author.Avatar != nil {
		avatarTmpl = m.Author.Avatar.Template(viewer)
	}
	if avatarTmpl.Title == "" {
		log.Printf("WARN message avatar title is empty: %v", m.Author.Avatar)
	}
	return template.MessageTemplate{
		IntermediateId:   utils.RandStringBytes(5),
		ChatId:           m.ChatId,
		MsgId:            m.Id,
		Quote:            &quoteTmpl,
		ViewerId:         viewer.Id,
		OwnerId:          owner.Id,
		AuthorId:         m.Author.Id,
		AuthorName:       m.Author.Name,
		AuthorAvatar:     avatarTmpl,
		Text:             m.Text,
		TextIntro:        utils.Shorten(m.Text, 69),
		MessageDropEvent: event.MessageDrop.FormatEventName(m.ChatId, m.Author.Id, m.Id),
	}, nil
}
