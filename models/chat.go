package models

import (
	"fmt"
	"sync"
)

type Chat struct {
	ID      int
	Name    string
	Owner   string
	users   []string
	history MessageStore
	mu      sync.Mutex
}

func (c *Chat) Log() string {
	return fmt.Sprintf("Chat{id:%d,name:[%s],owner:[%s]}", c.ID, c.Name, c.Owner)
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
		return nil, fmt.Errorf("only invited users can see users in chat")
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
		return nil, fmt.Errorf("only invited users can send messages")
	}
	msg := c.history.Add(message)
	return &msg, nil
}

func (c *Chat) GetMessage(user string, id int) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(user) {
		return nil, fmt.Errorf("only invited users can get messages")
	}
	return c.history.Get(id)
}

func (c *Chat) GetMessages(user string) ([]Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(user) {
		return nil, fmt.Errorf("only invited users can see messages")
	}
	return c.history.GetAll(), nil
}

func (c *Chat) RemoveMessage(user string, ID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isUserInChat(user) {
		return fmt.Errorf("only invited users can delete messages")
	}
	if !c.isAuthor(user, ID) && !c.isOwner(user) {
		return fmt.Errorf("only user that sent the original message or chat owner can delete messages")
	}
	return c.history.Delete(ID)
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
