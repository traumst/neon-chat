package template

import (
	"bytes"
	"html/template"
)

type HomeTemplate struct {
	Chats          []*ChatTemplate
	OpenTemplate   *ChatTemplate
	ActiveUser     string
	LoadLocal      bool
	ChatAddEvent   string
	IsAuthorized   bool
	LoginTemplate  LoginTemplate
	ChatCloseEvent string
}

func (h *HomeTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	homeTmpl := template.Must(template.ParseFiles(
		"static/html/home_page.html",
		"static/html/welcome_div.html",
		"static/html/user_settings_div.html",
		"static/html/nav/login_div.html",
		"static/html/nav/chat_li.html",
		"static/html/chat/chat_div.html",
		"static/html/chat/message_li.html"))
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
