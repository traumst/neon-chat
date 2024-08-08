package app

import (
	"sync"

	t "prplchat/src/model/template"
)

type Chat struct {
	Id        uint
	Name      string
	OwnerId   uint
	OwnerName string
	// users     []*User
	history MessageStore
	mu      sync.Mutex
}

// Parameters:
//
//	user: user who triggered state change
//	viewer: user who is viewing the chat
//	members: slice of users in chat
func (c *Chat) Template(user *User, viewer *User, members []*User) *t.ChatTemplate {
	var messages []t.MessageTemplate
	for _, msg := range c.history.GetAll() {
		if msg == nil {
			continue
		}
		messages = append(messages, *msg.Template(viewer))
	}
	users := make([]t.UserTemplate, len(members))
	for i, member := range members {
		users[i] = t.UserTemplate{
			ChatId:      c.Id,
			ChatOwnerId: c.OwnerId,
			UserId:      member.Id,
			UserName:    member.Name,
			UserEmail:   member.Email,
			ViewerId:    viewer.Id,
		}
	}
	usr := t.UserTemplate{
		ChatId:      c.Id,
		ChatOwnerId: c.OwnerId,
		UserId:      viewer.Id,
		UserName:    viewer.Name,
		ViewerId:    viewer.Id,
	}
	ownr := t.UserTemplate{
		ChatId:      c.Id,
		ChatOwnerId: c.OwnerId,
		UserId:      c.OwnerId,
		UserName:    c.OwnerName,
		ViewerId:    viewer.Id,
	}
	return &t.ChatTemplate{
		ChatId:   c.Id,
		ChatName: c.Name,
		User:     usr,
		Viewer:   usr,
		Owner:    ownr,
		Users:    users,
		Messages: messages,
	}
}
