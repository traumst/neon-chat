package template

import (
	"bytes"
	"fmt"
	"html/template"

	ti "prplchat/src/interface"
	"prplchat/src/model/event"
)

type HomeTemplate struct {
	IsAuthorized  bool
	Avatar        ti.Renderable
	User          ti.Renderable
	Chats         []ti.Renderable
	OpenChat      ti.Renderable
	LoginTemplate ti.Renderable
}

func (h HomeTemplate) ChatAddEvent() string {
	var openChatId uint = 0
	if h.OpenChat != nil {
		openChatId = h.OpenChat.(ChatTemplate).ChatId
	}
	return event.ChatAdd.FormatEventName(openChatId, h.User.(UserTemplate).UserId, 0)
}

func (h HomeTemplate) ChatInviteEvent() string {
	return string(event.ChatInvite)
}

func (h HomeTemplate) ChatCloseEvent() string {
	var openChatId uint = 0
	if h.OpenChat != nil {
		openChatId = h.OpenChat.(UserTemplate).ChatId
	}
	return event.ChatClose.FormatEventName(openChatId, h.User.(UserTemplate).UserId, 0)
}

func (h HomeTemplate) AvatarChangeEvent() string {
	return event.AvatarChange.FormatEventName(0, h.User.(UserTemplate).UserId, 0)
}

func (h HomeTemplate) HTML() (string, error) {
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
	if err := homeTmpl.Execute(&buf, h); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

// called from template
func (h HomeTemplate) ReverseChats() []ti.Renderable {
	len := len(h.Chats)
	chats := make([]ti.Renderable, len)
	for i := 0; i < len; i += 1 {
		chats[i] = h.Chats[(len-1)-i]
	}
	return chats
}
