package model

import (
	"fmt"
	"strings"
	"sync"
)

type MessageStore struct {
	mu       sync.Mutex
	messages []*Message
	nextID   int
}

func (store *MessageStore) Add(message *Message) (*Message, error) {
	author := strings.TrimSpace(message.Author)
	msg := strings.TrimSpace(message.Text)
	if author == "" || msg == "" {
		return nil, fmt.Errorf("bad arguments")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	message.ID = store.nextID
	store.messages = append(store.messages, message)
	store.nextID += 1
	return message, nil
}

func (store *MessageStore) GetAll() []*Message {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.messages
}

func (store *MessageStore) Get(id int) (*Message, error) {
	if id < 0 || id >= len(store.messages) {
		return nil, fmt.Errorf("message not found")
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	return store.messages[id], nil
}

func (store *MessageStore) Delete(id int) error {
	if id < 0 || id >= len(store.messages) {
		return fmt.Errorf("message not found")
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	store.messages = append(store.messages[id:], store.messages[id+1:]...)
	return nil
}
