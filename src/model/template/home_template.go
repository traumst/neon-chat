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
}

func (h *HomeTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	homeTmpl := template.Must(template.ParseFiles(
		"html/home_page.html",
		"html/bits/welcome_div.html",
		"html/bits/chat_div.html",
		"html/bits/chat_li.html"))
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
