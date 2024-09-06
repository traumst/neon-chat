package app

import (
	"fmt"
	t "neon-chat/src/model/template"
	"neon-chat/src/utils"
)

type Quote struct {
	Id     uint
	ChatId uint
	Author *User
	Text   string
	Quote  *Message
}

func (m *Quote) Template(
	viewer *User,
	owner *User,
	avatar *Avatar,
) (t.QuoteTemplate, error) {
	if viewer == nil || viewer.Id == 0 {
		return t.QuoteTemplate{}, fmt.Errorf("viewer cannot be nil or blank")
	}
	if m.Author == nil || m.Author.Id == 0 || m.Author.Name == "" {
		return t.QuoteTemplate{}, fmt.Errorf("message author cannot be nil or blank")
	}
	if m.ChatId == 0 {
		return t.QuoteTemplate{}, fmt.Errorf("message chatId cannot be 0")
	}
	if m.Id == 0 {
		return t.QuoteTemplate{}, fmt.Errorf("message chatId and Id cannot be 0")
	}
	return t.QuoteTemplate{
		IntermediateId: utils.RandStringBytes(5),
		ChatId:         m.ChatId,
		MsgId:          m.Id,
		AuthorId:       m.Author.Id,
		AuthorName:     m.Author.Name,
		AuthorAvatar:   avatar.Template(viewer),
		Text:           m.Text,
		TextIntro:      utils.Shorten(m.Text, 69),
	}, nil
}
