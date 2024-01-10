package models

import (
	"sync"
)

type ChatCollapsed struct {
	ID    int
	Name  string
	Owner string
}

type Chat struct {
	ID      int
	Name    string
	Owner   string
	users   []string
	history MessageStore
	mu      sync.Mutex
}

func (c *Chat) AddUser(user string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users = append(c.users, user)
}

func (c *Chat) RemoveUser(user string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, u := range c.users {
		if u == user {
			c.users = append(c.users[:i], c.users[i+1:]...)
			break
		}
	}
}

func (c *Chat) AddMessage(message Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.history.Add(message)
}

func (c *Chat) GetMessages() []Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.history.Get()
}

func (c *Chat) RemoveMessage(ID int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.history.Delete(ID)
}
