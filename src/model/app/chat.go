package app

import (
	"fmt"
	"sync"

	t "go.chat/src/model/template"
)

type Chat struct {
	Id      int
	Name    string
	Owner   *User
	users   []*User
	history MessageStore
	mu      sync.Mutex
}

// type ChatTable struct {
// 	Id      uint   `db:"id"`
// 	Name    string `db:"name"`
// 	OwnerId uint   `db:"ownerId"`
// }

// type ChatUsersTable struct {
// 	ChatId int  `db:"chatId"`
// 	UserId uint `db:"userId"`
// }

func (c *Chat) isOwner(userId uint) bool {
	return c.Owner.Id == userId
}

func (c *Chat) isAuthor(userId uint, msgId int) bool {
	msg, _ := c.history.Get(msgId)
	if msg != nil && msg.Id == msgId {
		return msg.Author.Id == userId
	}
	return false
}

func (c *Chat) isUserInChat(userId uint) bool {
	for _, u := range c.users {
		if u.Id == userId {
			return true
		}
	}
	return false
}

func (c *Chat) AddUser(ownerId uint, user *User) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isOwner(ownerId) {
		return fmt.Errorf("only the owner can invite users")
	}
	if c.isUserInChat(user.Id) {
		return fmt.Errorf("user already in chat")
	}
	c.users = append(c.users, user)
	return nil
}

func (c *Chat) SyncUser(user *User) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, u := range c.users {
		if u.Id == user.Id {
			u.Name = user.Name
			return nil
		}
	}
	return fmt.Errorf("user[%d] is not in chat[%d]", user.Id, c.Id)
}

func (c *Chat) GetUsers(userId uint) ([]*User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(userId) {
		return nil, fmt.Errorf("user[%d] is not in chat[%d]", userId, c.Id)
	}
	return c.users, nil
}

func (c *Chat) RemoveUser(ownerId uint, userId uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isOwner(ownerId) && ownerId != userId {
		return fmt.Errorf("only the owner can remove users from chat")
	}
	if !c.isUserInChat(userId) {
		return fmt.Errorf("only invited users can be removed from chat")
	}
	for i, u := range c.users {
		if u.Id == userId {
			c.users = append(c.users[:i], c.users[i+1:]...)
			break
		}
	}
	return nil
}

func (c *Chat) AddMessage(userId uint, message Message) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(userId) {
		return nil, fmt.Errorf("only invited users can add messages")
	}

	return c.history.Add(&message)
}

func (c *Chat) GetMessage(userId uint, msgId int) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(userId) {
		return nil, fmt.Errorf("only invited users can get messages")
	}
	return c.history.Get(msgId)
}

func (c *Chat) DropMessage(userId uint, msgId int) (*Message, error) {
	msg, err := c.GetMessage(userId, msgId)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(userId) {
		return msg, fmt.Errorf("only invited users can delete messages")
	}
	if !c.isAuthor(userId, msgId) && !c.isOwner(userId) {
		return msg, fmt.Errorf("only user that sent the original message or chat owner can delete messages")
	}
	return msg, c.history.Delete(msg)
}

func (c *Chat) Template(user *User, viewer *User) *t.ChatTemplate {
	var messages []t.MessageTemplate
	for _, msg := range c.history.GetAll() {
		if msg == nil {
			continue
		}
		messages = append(messages, *msg.Template(user))
	}
	users := make([]t.UserTemplate, len(c.users))
	for i, u := range c.users {
		users[i] = t.UserTemplate{
			ChatId:      c.Id,
			ChatOwnerId: c.Owner.Id,
			UserId:      u.Id,
			UserName:    u.Name,
			ViewerId:    viewer.Id,
		}
	}
	usr := t.UserTemplate{
		ChatId:      c.Id,
		ChatOwnerId: c.Owner.Id,
		UserId:      user.Id,
		UserName:    user.Name,
		ViewerId:    viewer.Id,
	}
	ownr := t.UserTemplate{
		ChatId:      c.Id,
		ChatOwnerId: c.Owner.Id,
		UserId:      c.Owner.Id,
		UserName:    c.Owner.Name,
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
