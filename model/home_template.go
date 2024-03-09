package model

import (
	"bytes"
	"html/template"
)

type HomeTemplate struct {
	Chats        []*ChatTemplate
	OpenTemplate *ChatTemplate
	ActiveUser   string
}

func (h *HomeTemplate) GetHTML() (string, error) {
	var buf bytes.Buffer
	homeTmpl := template.Must(template.ParseFiles(
		"html/home.html",
		"html/welcome.html",
		"html/chat_div.html",
		"html/chat_li.html"))
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
