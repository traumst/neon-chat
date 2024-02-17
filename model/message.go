package model

type Message struct {
	ID     int
	Author string
	Text   string
}

func (c *Message) ToTemplate(user string) *MessageTemplate {
	return &MessageTemplate{
		ID:         c.ID,
		Author:     c.Author,
		Text:       c.Text,
		ActiveUser: user,
	}
}
