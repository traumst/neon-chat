package model

import (
	"bytes"
	"html/template"
)

type ChatTemplate struct {
	ChatID   int
	Name     string
	User     string
	Owner    string
	Users    []string
	Messages []MessageTemplate
}

func (c *ChatTemplate) GetHTML() (string, error) {
	var buf bytes.Buffer
	chatTmpl := template.Must(template.ParseFiles("html/chat.html", "html/message.html"))
	err := chatTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *ChatTemplate) GetShortHTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("html/chat_li.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
