package app

import (
	"fmt"
	"sync"

	"go.chat/model/template"
)

type Chat struct {
	ID      int
	Name    string
	Owner   string
	users   []string
	history MessageStore
	mu      sync.Mutex
}

func (c *Chat) isOwner(user string) bool {
	return user == c.Owner
}

func (c *Chat) isAuthor(user string, msgID int) bool {
	msg, _ := c.history.Get(msgID)
	if msg != nil && msg.ID == msgID {
		return msg.Author == user
	}
	return false
}

func (c *Chat) isUserInChat(user string) bool {
	for _, u := range c.users {
		if u == user {
			return true
		}
	}
	return false
}

func (c *Chat) AddUser(owner string, user string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isOwner(owner) {
		return fmt.Errorf("only the owner can invite users")
	}
	c.users = append(c.users, user)
	return nil
}

func (c *Chat) GetUsers(user string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(user) {
		return nil, fmt.Errorf("only invited users can see users in chat, %s", user)
	}
	return c.users, nil
}

func (c *Chat) RemoveUser(owner string, user string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isOwner(owner) {
		return fmt.Errorf("only the owner can remove users from chat")
	}
	if !c.isUserInChat(user) {
		return fmt.Errorf("only invited users can be removed from chat")
	}
	for i, u := range c.users {
		if u == user {
			c.users = append(c.users[:i], c.users[i+1:]...)
			break
		}
	}
	return nil
}

func (c *Chat) AddMessage(user string, message Message) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(user) {
		return nil, fmt.Errorf("only invited users can add messages")
	}

	return c.history.Add(&message)
}

func (c *Chat) GetMessage(user string, id int) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(user) {
		return nil, fmt.Errorf("only invited users can get messages")
	}
	return c.history.Get(id)
}

func (c *Chat) DropMessage(user string, ID int) (*Message, error) {
	msg, err := c.GetMessage(user, ID)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(user) {
		return msg, fmt.Errorf("only invited users can delete messages")
	}
	if !c.isAuthor(user, ID) && !c.isOwner(user) {
		return msg, fmt.Errorf("only user that sent the original message or chat owner can delete messages")
	}
	return msg, c.history.Delete(msg)
}

func (c *Chat) Template(user string) *template.ChatTemplate {
	var messages []template.MessageTemplate
	for _, msg := range c.history.GetAll() {
		if msg == nil {
			continue
		}
		messages = append(messages, *msg.Template(user))
	}
	return &template.ChatTemplate{
		ChatID:   c.ID,
		Name:     c.Name,
		User:     user,
		Viewer:   user,
		Owner:    c.Owner,
		Users:    c.users[0:],
		Messages: messages,
	}
}

func (c *Chat) String() string {
	return fmt.Sprintf("Chat{ID:%d,Name:%s,Owner:%s", c.ID, c.Name, c.Owner)
}
