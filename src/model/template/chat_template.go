package template

import (
	"bytes"
	"html/template"
)

type ChatTemplate struct {
	ChatId   int
	Name     string
	User     UserTemplate
	Viewer   UserTemplate
	Owner    UserTemplate
	Users    []UserTemplate
	Messages []MessageTemplate
}

func (c *ChatTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	chatTmpl := template.Must(template.ParseFiles(
		"static/html/bits/chat_div.html",
		"static/html/bits/members_div.html",
		"static/html/bits/message_li.html"))
	err := chatTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *ChatTemplate) ShortHTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/bits/chat_li.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
