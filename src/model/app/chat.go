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

// Short html template does not require members or messages
//
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
	// current viewer + chat owner
	var usr t.UserTemplate
	var ownr t.UserTemplate
	if viewer == nil {
		log.Printf("Chat.Template ERROR viewer cannot be nil\n")
		return nil
	}
	usr = t.UserTemplate{
		ChatId:      c.Id,
		ChatOwnerId: c.OwnerId,
		UserId:      viewer.Id,
		UserName:    viewer.Name,
		UserEmail:   viewer.Email,
		ViewerId:    viewer.Id,
	}
	ownr = t.UserTemplate{
		ChatId:      c.Id,
		ChatOwnerId: c.OwnerId,
		UserId:      c.OwnerId,
		UserName:    c.OwnerName,
		//UserEmail:   c.OwnerEmail,
		ViewerId: viewer.Id,
	}
	// chat users
	users := make([]t.UserTemplate, 0)
	if len(members) <= 0 {
		log.Printf("Chat.Template INFO chat[%d] has no users\n", c.Id)
	} else {
		for _, member := range members {
			if member == nil {
				log.Printf("Chat.Template TRACE skip nil member in chat[%d]\n", c.Id)
				continue
			}
			users = append(users, t.UserTemplate{
				ChatId:      c.Id,
				ChatOwnerId: c.OwnerId,
				UserId:      member.Id,
				UserName:    member.Name,
				UserEmail:   member.Email,
				ViewerId:    viewer.Id,
			})
		}
	}
	// chat messages
	messages := make([]t.MessageTemplate, 0)
	if len(msgs) > 0 {
		for idx, msg := range msgs {
			if msg == nil {
				log.Printf("Chat.Template TRACE skip nil msg on index[%d] in chat[%d]\n", idx, c.Id)
				continue
			}
			msgTmpl, err := msg.Template(viewer, &User{Id: c.OwnerId, Name: c.OwnerName})
			if err != nil {
				log.Printf("Chat.Template ERROR failed to create message template, %s\n", err)
				continue
			}
			messages = append(messages, *msgTmpl)
		}
	} else {
		log.Printf("Chat.Template INFO chat[%d] has no messages\n", c.Id)
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
