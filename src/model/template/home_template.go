package template

import (
	"bytes"
	"fmt"
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
	var openChatId uint = 0
	if h.OpenChat != nil {
		openChatId = h.OpenChat.ChatId
	}
	return event.ChatAdd.FormatEventName(openChatId, h.User.UserId, 0)
}

func (h *HomeTemplate) ChatInviteEvent() string {
	return string(event.ChatInvite)
}

func (h *HomeTemplate) ChatCloseEvent() string {
	var openChatId uint = 0
	if h.OpenChat != nil {
		openChatId = h.OpenChat.ChatId
	}
	return event.ChatClose.FormatEventName(openChatId, h.User.UserId, 0)
}

func (h *HomeTemplate) AvatarChangeEvent() string {
	return event.AvatarChange.FormatEventName(0, h.User.UserId, 0)
}

func (h *HomeTemplate) HTML() (string, error) {
	if err := h.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
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
func (h *HomeTemplate) ReverseChats() []*ChatTemplate {
	len := len(h.Chats)
	chats := make([]*ChatTemplate, len)
	for i := 0; i < len; i += 1 {
		chats[i] = h.Chats[(len-1)-i]
	}
	return chats
}

func (h *HomeTemplate) validate() error {
	if h.Avatar == nil {
		return fmt.Errorf("HomeTemplate requires Avatar, but is nil")
	}
	if h.User.UserId < 1 {
		return fmt.Errorf("HomeTemplate requires User.UserId, but is [%d]", h.User.UserId)
	}
	if h.User.UserName == "" {
		return fmt.Errorf("HomeTemplate requires User.UserName, but is [%s]", h.User.UserName)
	}
	if len(h.Chats) < 1 {
		return fmt.Errorf("HomeTemplate requires Chats, but is empty")
	}
	if h.OpenChat == nil {
		return fmt.Errorf("HomeTemplate requires OpenChat, but is nil")
	}
	return nil
}
