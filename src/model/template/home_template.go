package template

import (
	"bytes"
	"html/template"

	"prplchat/src/model/event"
)

type HomeTemplate struct {
	IsAuthorized  bool
	Avatar        *AvatarTemplate
	User          UserTemplate
	Chats         []*ChatTemplate
	OpenChat      *ChatTemplate
	LoginTemplate AuthTemplate
}

func (h *HomeTemplate) ChatAddEvent() string {
	var openChatId int = -1
	if h.OpenChat != nil {
		openChatId = h.OpenChat.ChatId
	}
	return event.ChatAdd.FormatEventName(openChatId, h.User.UserId, -5)
}

func (h *HomeTemplate) ChatInviteEvent() string {
	return string(event.ChatInvite)
}

func (h *HomeTemplate) ChatCloseEvent() string {
	var openChatId int = -1
	if h.OpenChat != nil {
		openChatId = h.OpenChat.ChatId
	}
	return event.ChatClose.FormatEventName(openChatId, h.User.UserId, -6)
}

func (h *HomeTemplate) AvatarChangeEvent() string {
	return event.AvatarChange.FormatEventName(0, h.User.UserId, -5)
}

func (h *HomeTemplate) HTML() (string, error) {
	homeTmpl := template.Must(template.ParseFiles(
		"static/html/home_page.html",
		// left panel
		"static/html/left_panel.html",
		"static/html/avatar_div.html",
		"static/html/nav/auth_div.html",
		"static/html/nav/chat_li.html",
		// right panel
		"static/html/utils/user_settings_div.html",
		"static/html/chat/welcome_div.html",
		"static/html/chat/chat_div.html",
		"static/html/user_div.html",
		"static/html/chat/message_li.html",
	))
	var buf bytes.Buffer
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
