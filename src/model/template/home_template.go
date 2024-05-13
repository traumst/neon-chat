package template

import (
	"bytes"
	"html/template"

	"go.chat/src/model/event"
)

type HomeTemplate struct {
	Chats         []*ChatTemplate
	OpenChat      *ChatTemplate
	User          UserTemplate
	LoadLocal     bool
	IsAuthorized  bool
	LoginTemplate LoginTemplate
	Avatar        *AvatarTemplate
}

func (h *HomeTemplate) ChatAddEvent() string {
	var openChatId int = -1
	if h.OpenChat != nil {
		openChatId = h.OpenChat.ChatId
	}
	return event.ChatAdd.FormatEventName(openChatId, h.User.UserId, -5)
}

func (h *HomeTemplate) ChatCloseEvent() string {
	var openChatId int = -1
	if h.OpenChat != nil {
		openChatId = h.OpenChat.ChatId
	}
	return event.ChatClose.FormatEventName(openChatId, h.User.UserId, -6)
}

func (h *HomeTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	homeTmpl := template.Must(template.ParseFiles(
		"static/html/home_page.html",
		// left panel
		"static/html/left_panel.html",
		"static/html/avatar_div.html",
		"static/html/utils/user_settings_div.html",
		"static/html/nav/login_div.html",
		"static/html/nav/chat_li.html",
		// right panel
		"static/html/welcome_div.html",
		"static/html/chat/chat_div.html",
		"static/html/user_div.html",
		"static/html/chat/message_li.html",
	))
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
