package app

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
	author := strings.TrimSpace(m.Author.Name)
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
	return s.nonNil()
}

func (s *MessageStore) Get(id int) (*Message, error) {
	if id < 0 || id >= len(s.messages) {
		return nil, fmt.Errorf("invalid msg id")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	msg := s.messages[id]
	if msg == nil {
		return nil, fmt.Errorf("msg not found")
	}

	return msg, nil
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
	s.messages[m.ID] = nil
	return nil
}

func (s *MessageStore) nonNil() []*Message {
	var nonNil []*Message
	for _, msg := range s.messages {
		if msg != nil {
			nonNil = append(nonNil, msg)
		}
	}
	return nonNil
}
