package template

import (
	"bytes"
	"html/template"
)

type HomeTemplate struct {
	Chats        []*ChatTemplate
	OpenTemplate *ChatTemplate
	ActiveUser   string
	LoadLocal    bool
	ChatAddEvent string
}

func (h *HomeTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	homeTmpl := template.Must(template.ParseFiles(
		"static/html/home_page.html",
		"static/html/bits/welcome_div.html",
		"static/html/bits/chat_div.html",
		"static/html/bits/chat_li.html"))
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
