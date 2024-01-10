package models

import "sync"

type MessageStore struct {
	mu       sync.Mutex
	messages []Message
	nextID   int
}

func (store *MessageStore) Add(message Message) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.nextID += 1
	message.ID = store.nextID
	store.messages = append(store.messages, message)
}

func (store *MessageStore) Get() []Message {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.messages
}

func (store *MessageStore) Delete(id int) {
	store.mu.Lock()
	defer store.mu.Unlock()
	for i, message := range store.messages {
		if message.ID == id {
			store.messages = append(store.messages[:i], store.messages[i+1:]...)
			break
		}
	}
}