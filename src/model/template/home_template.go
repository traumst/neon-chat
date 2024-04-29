package template

import (
	"bytes"
	"html/template"
)

type HomeTemplate struct {
	Chats         []*ChatTemplate
	OpenTemplate  *ChatTemplate
	ActiveUser    string
	LoadLocal     bool
	ChatAddEvent  string
	IsAuthorized  bool
	LoginTemplate LoginTemplate
}

func (h *HomeTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	homeTmpl := template.Must(template.ParseFiles(
		"static/html/home_page.html",
		"static/html/bits/welcome_div.html",
		"static/html/bits/login_div.html",
		"static/html/bits/chat_div.html",
		"static/html/bits/chat_li.html",
		"static/html/bits/message_li.html"))
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// called from template
func (h *HomeTemplate) ReverseChats() []*ChatTemplate {
	len := len(h.Chats)
	chats := make([]*ChatTemplate, len)
	for i := 0; i < len; i += 1 {
		chats[i] = h.Chats[(len-1)-i]
	}
	return chats
}
