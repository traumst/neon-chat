package model

import (
	"fmt"
)

type Message struct {
	ID     int
	Author string
	Text   string
}

func (m *Message) Log() string {
	if m == nil {
		return "Message{nil}"
	}
	return fmt.Sprintf("Message{id:%d,author:[%s],text:[%s]}", m.ID, m.Author, m.Text)
}

func (c *Message) ToTemplate(user string) *MessageTemplate {
	return &MessageTemplate{
		ID:         c.ID,
		Author:     c.Author,
		Text:       c.Text,
		ActiveUser: user,
	}
}
