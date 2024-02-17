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

func (s *MessageStore) Add(m *Message) (*Message, error) {
	author := strings.TrimSpace(m.Author)
	msg := strings.TrimSpace(m.Text)
	if author == "" || msg == "" {
		return nil, fmt.Errorf("bad arguments")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	m.ID = s.nextID
	s.messages = append(s.messages, m)
	s.nextID += 1
	return m, nil
}

func (s *MessageStore) GetAll() []*Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.messages
}

func (s *MessageStore) Get(id int) (*Message, error) {
	if id < 0 || id >= len(s.messages) {
		return nil, fmt.Errorf("message not found")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.messages[id], nil
}

func (s *MessageStore) Delete(m *Message) error {
	if m == nil {
		return fmt.Errorf("cannot remove NIL")
	}
	if m.ID < 0 || m.ID >= len(s.messages) {
		return fmt.Errorf("message not found")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages[m.ID:], s.messages[m.ID+1:]...)
	return nil
}
