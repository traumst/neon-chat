package template

import (
	"bytes"
	"html/template"
)

type ChatTemplate struct {
	ChatID   int
	Name     string
	User     string
	Viewer   string
	Owner    string
	Users    []string
	Messages []MessageTemplate
}

func (c *ChatTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	chatTmpl := template.Must(template.ParseFiles(
		"html/bits/chat_div.html",
		"html/bits/members_div.html",
		"html/bits/member_div.html",
		"html/bits/message_li.html"))
	err := chatTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *ChatTemplate) ShortHTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("html/bits/chat_li.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}