package app

import (
	"log"
	t "prplchat/src/model/template"
)

type Chat struct {
	Id        uint
	Name      string
	OwnerId   uint
	OwnerName string
}

// Parameters:
//
//	user: user who triggered state change
//	viewer: user who is viewing the chat
//	members: slice of users in chat
func (c *Chat) Template(
	user *User,
	viewer *User,
	members []*User,
	msgs []*Message,
) *t.ChatTemplate {
	// chat messages
	msgCount := len(msgs)
	messages := make([]t.MessageTemplate, msgCount)
	if msgCount > 0 {
		for _, msg := range msgs {
			if msg == nil {
				continue
			}
			messages = append(messages, *msg.Template(viewer))
		}
	} else {
		log.Printf("Chat.Template INFO chat[%d] has no messages\n", c.Id)
	}
	// chat users
	userCount := len(members)
	users := make([]t.UserTemplate, userCount)
	if userCount > 0 {
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
	} else {
		log.Printf("Chat.Template INFO chat[%d] has no users\n", c.Id)
		return nil
	}
	// current viewer + chat owner
	var usr t.UserTemplate
	var ownr t.UserTemplate
	if viewer != nil {
		usr = t.UserTemplate{
			ChatId:      c.Id,
			ChatOwnerId: c.OwnerId,
			UserId:      viewer.Id,
			UserName:    viewer.Name,
			ViewerId:    viewer.Id,
		}
		ownr = t.UserTemplate{
			ChatId:      c.Id,
			ChatOwnerId: c.OwnerId,
			UserId:      c.OwnerId,
			UserName:    c.OwnerName,
			ViewerId:    viewer.Id,
		}
	} else {
		log.Printf("Chat.Template ERROR viewer cannot be nil\n")
		return nil
	}
	// chat
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
