package app

import (
	"neon-chat/src/interfaces"
	"neon-chat/src/model/template"
)

type Home struct {
	IsAuthorized bool
	User         *User
	Avatar       *Avatar
	UserChats    []*Chat
	OpenChat     *Chat
}

func (h *Home) Template() template.HomeTemplate {
	var chatTemplates []interfaces.Renderable
	for _, chat := range h.UserChats {
		chatTemplates = append(chatTemplates, chat.Template(h.User, h.User, nil, nil))
	}
	var openChatId uint
	var chatOwnerId uint
	if h.OpenChat != nil {
		openChatId = h.OpenChat.Id
		chatOwnerId = h.OpenChat.OwnerId
	}
	userTemplate := h.User.Template(openChatId, chatOwnerId, h.User.Id)
	return template.HomeTemplate{
		Chats:         chatTemplates,
		OpenChat:      h.OpenChat.Template(h.User, h.User, nil, nil),
		User:          userTemplate,
		IsAuthorized:  h.IsAuthorized,
		LoginTemplate: template.AuthTemplate{},
		Avatar:        h.Avatar.Template(h.User),
	}
}
