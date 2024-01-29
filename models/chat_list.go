package models

import (
	"fmt"
	"sync"

	"go.chat/utils"
)

type ChatList struct {
	mu     sync.Mutex
	chats  []*Chat
	userAt map[string]*Chat
	nextID int
	isInit bool
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
	cl.init(owner)
	chat := Chat{
		ID:      cl.nextID,
		Name:    chatName,
		Owner:   owner,
		users:   []string{owner},
		history: MessageStore{},
	}

	cl.chats = append(cl.chats, &chat)
	cl.nextID += 1
	cl.userAt[owner] = cl.chats[chat.ID]
	return chat.ID
}

func (cl *ChatList) OpenChat(user string, id int) (*Chat, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if id < 0 || id >= len(cl.chats) {
		return nil, fmt.Errorf("invalid chat id")
	}
	cl.init(user)
	openChat := cl.chats[id]
	if !openChat.isOwner(user) && !openChat.isUserInChat(user) {
		return nil, fmt.Errorf("user[%s] is not in chat[%d]", user, id)
	}
	cl.userAt[user] = cl.chats[id]
	return openChat, nil
}

func (cl *ChatList) GetOpenChat(user string) *Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init(user)
	return cl.userAt[user]
}

func (cl *ChatList) GetChats(user string) []*Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	var userChats []*Chat
	for _, chat := range cl.chats {
		if chat.isOwner(user) || chat.isUserInChat(user) {
			userChats = append(userChats, chat)
		}
	}
	return userChats
}

func (cl *ChatList) GetExcept(user string, chatIDs []int) []*Chat {
	var userChats []*Chat
	for _, chat := range cl.GetChats(user) {
		if !chat.isOwner(user) && !chat.isUserInChat(user) {
			continue
		}
		if utils.Contains(chatIDs, chat.ID) {
			continue
		}
		userChats = append(userChats, chat)
	}
	return userChats
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

func (cl *ChatList) InviteUser(user string, chatID int, invitee string) (string, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init(user)
	cl.init(invitee)
	if chatID < 0 || chatID >= len(cl.chats) {
		return "", fmt.Errorf("invalid chat index[%d]", chatID)
	}
	chat := cl.chats[chatID]
	if !chat.isOwner(user) {
		return "", fmt.Errorf("user[%s] is not owner of chat %d", user, chatID)
	}
	chat.users = append(chat.users, invitee)
	return invitee, nil
}

func (cl *ChatList) init(user string) {
	if cl.isInit {
		return
	}

	cl.chats = []*Chat{}
	cl.userAt = make(map[string]*Chat)
	cl.isInit = true
}
