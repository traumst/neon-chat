package models

import (
	"fmt"
	"sync"
)

type MessageStore struct {
	mu       sync.Mutex
	messages []Message
	nextID   int
}

func (store *MessageStore) Add(message Message) Message {
	store.mu.Lock()
	defer store.mu.Unlock()
	message.ID = store.nextID
	store.messages = append(store.messages, message)
	store.nextID += 1
	return message
}

func (store *MessageStore) GetAll() []Message {
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
	return &store.messages[id], nil
}

func (store *MessageStore) Delete(id int) error {
	if id < 0 || id >= len(store.messages) {
		return fmt.Errorf("message not found")
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	for i, message := range store.messages {
		if message.ID == id {
			store.messages = append(store.messages[:i], store.messages[i+1:]...)
			break
		}
	}
	return nil
}
