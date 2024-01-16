package models

import (
	"sync"
)

type ChatList struct {
	mu     sync.Mutex
	chats  []*Chat
	open   *Chat
	nextID int
}

func (h *ChatList) AddChat(owner string, name string) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	chat := Chat{
		ID:      h.nextID,
		Name:    name,
		Owner:   owner,
		users:   []string{owner},
		history: MessageStore{},
	}
	h.chats = append(h.chats, &chat)
	h.nextID += 1
	return chat.ID
}

func (h *ChatList) OpenChat(id int) *Chat {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.open = h.chats[id]
	return h.open
}

func (h *ChatList) GetOpenChat() *Chat {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.open
}

func (h *ChatList) GetChats() []*Chat {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.chats
}

func (h *ChatList) DeleteChat(index int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.chats = append(h.chats[:index], h.chats[index+1:]...)
}
