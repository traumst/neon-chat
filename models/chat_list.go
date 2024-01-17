package models

import (
	"fmt"
	"sync"
)

type ChatList struct {
	mu        sync.Mutex
	chats     []*Chat
	userAt    map[string]*Chat
	userChats map[string][]*Chat
	nextID    int
}

type ChatError struct {
	msg string
}

func (e *ChatError) Error() string {
	return e.msg
}

func (cl *ChatList) AddChat(owner string, chatName string) int {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	chat := Chat{
		ID:      cl.nextID,
		Name:    chatName,
		Owner:   owner,
		users:   []string{owner},
		history: MessageStore{},
	}
	cl.chats = append(cl.chats, &chat)
	if cl.userAt == nil {
		cl.userAt = make(map[string]*Chat)
	}
	if cl.userChats == nil {
		cl.userChats = make(map[string][]*Chat)
	}
	if cl.userChats[owner] == nil {
		cl.userChats[owner] = make([]*Chat, 0)
	}
	cl.userChats[owner] = append(cl.userChats[owner], &chat)
	cl.nextID += 1
	cl.userAt[owner] = cl.chats[chat.ID]
	return chat.ID
}

func (cl *ChatList) OpenChat(user string, id int) (*ChatTemplate, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if id < 0 || id >= len(cl.chats) {
		return nil, fmt.Errorf("invalid chat id")
	}
	openChat := cl.userAt[user]
	if !openChat.isOwner(user) && !openChat.isUserInChat(user) {
		return nil, fmt.Errorf("user[%s] is not in chat[%d]", user, id)
	}
	cl.userAt[user] = cl.chats[id]
	return &ChatTemplate{
		Chat:       cl.userAt[user],
		ActiveUser: user,
	}, nil
}

func (cl *ChatList) OpenTemplate(user string) (*ChatTemplate, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	openChat := cl.userAt[user]
	if openChat == nil {
		return nil, fmt.Errorf("user[%s] has no open chat", user)
	}
	return &ChatTemplate{
		Chat:       openChat,
		ActiveUser: user,
	}, nil
}

func (cl *ChatList) GetChats(user string) []*Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.userChats[user]
}

func (cl *ChatList) DeleteChat(user string, index int) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if index < 0 || index >= len(cl.chats) {
		return fmt.Errorf("invalid chat index[%d]", index)
	}
	removed := cl.chats[index]
	if !removed.isOwner(user) {
		return fmt.Errorf("user[%s] is not owner of chat %d", user, index)
	}
	cl.chats = append(cl.chats[:index], cl.chats[index+1:]...)
	return nil
}
