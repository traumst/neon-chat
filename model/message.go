package model

type Message struct {
	ID     int
	ChatID int
	Author string
	Text   string
}

func (m *Message) ToTemplate(user string) *MessageTemplate {
	return &MessageTemplate{
		MsgID:      m.ID,
		ChatID:     m.ChatID,
		Author:     m.Author,
		Text:       m.Text,
		ActiveUser: user,
	}
}
