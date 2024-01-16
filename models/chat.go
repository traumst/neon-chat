package models

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"text/template"
)

type Chat struct {
	ID      int
	Name    string
	Owner   string
	users   []string
	history MessageStore
	mu      sync.Mutex
}

var chatTmpl = template.Must(template.ParseFiles("views/chat.html"))

func (c *Chat) GetHTML() (string, error) {
	var buf bytes.Buffer
	err := chatTmpl.Execute(&buf, c)
	if err != nil {
		log.Printf("------ GetHTML ERROR template, %s\n", c.Log())
		return "", err
	}
	return buf.String(), nil
}

func (c *Chat) Log() string {
	return fmt.Sprintf("Chat{id:%d,name:[%s],owner:[%s]}", c.ID, c.Name, c.Owner)
}

func (c *Chat) AddUser(user string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users = append(c.users, user)
}

func (c *Chat) GetUsers() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.users
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

func (c *Chat) AddMessage(message Message) Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.history.Add(message)
}

func (c *Chat) GetMessage(id int) Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.history.Get()[id]
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
