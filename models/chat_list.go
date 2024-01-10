package models

import "sync"

type ChatList struct {
	mu     sync.Mutex
	chats  []*Chat
	open   *Chat
	nextID int
}

func (h *ChatList) AddChat(owner string, name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.nextID += 1
	chat := Chat{
		ID:      h.nextID,
		Name:    name,
		Owner:   owner,
		users:   []string{owner},
		history: MessageStore{},
	}
	h.chats = append(h.chats, &chat)
}

func (h *ChatList) GetOpenChat() *Chat {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.open
}

func (h *ChatList) GetChatsCollapsed() []ChatCollapsed {
	h.mu.Lock()
	defer h.mu.Unlock()

	collapsed := []ChatCollapsed{}
	for _, chat := range h.chats {
		collapsed = append(collapsed, ChatCollapsed{
			ID:    chat.ID,
			Name:  chat.Name,
			Owner: chat.Owner,
		})
	}
	return collapsed
}

func (h *ChatList) DeleteChat(index int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.chats = append(h.chats[:index], h.chats[index+1:]...)
}
