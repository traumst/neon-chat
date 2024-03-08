package model

type Message struct {
	ID     int
	Author string
	Text   string
}

func (m *Message) ToTemplate(user string) *MessageTemplate {
	return &MessageTemplate{
		MsgID:      m.ID,
		Author:     m.Author,
		Text:       m.Text,
		ActiveUser: user,
	}
}
