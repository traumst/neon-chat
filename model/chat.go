package model

import (
	"fmt"
	"log"
	"strings"
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
	log.Printf("------ Chat.isUserInChat TRACE user[%s], chat[%s]\n", user, strings.Join(c.users, ","))
	for _, u := range c.users {
		if u == user {
			return true
		}
	}
	return false
}

func (c *Chat) Log() string {
	if c == nil {
		return "Chat: NIL"
	}
	return fmt.Sprintf("Chat{id:%d,name:[%s],owner:[%s]}", c.ID, c.Name, c.Owner)
}

func (c *Chat) ToTemplate(user string) *ChatTemplate {
	messages := make([]MessageTemplate, 0)
	for _, msg := range c.history.GetAll() {
		messages = append(messages, *msg.ToTemplate(user))
	}
	return &ChatTemplate{
		ID:         c.ID,
		Name:       c.Name,
		ActiveUser: user,
		Users:      c.users,
		Messages:   messages,
	}
}

func (c *Chat) AddUser(owner string, user string) error {
	log.Printf("------ Chat.AddUser TRACE user[%s] added by owner[%s] to chat[%s]\n",
		user, owner, strings.Join(c.users, ","))
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
