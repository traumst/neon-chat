package model

import (
	"fmt"
	"sync"
)

type ChatList struct {
	mu     sync.Mutex
	chats  []*Chat
	open   map[string]*Chat
	nextID int
	isInit bool
}

func (cl *ChatList) init() {
	if cl.isInit {
		return
	}

	cl.isInit = true
	cl.chats = []*Chat{}
	cl.open = make(map[string]*Chat)
}

func (cl *ChatList) AddChat(owner string, chatName string) int {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	chat := Chat{
		ID:      cl.nextID,
		Name:    chatName,
		Owner:   owner,
		users:   []string{owner},
		history: MessageStore{},
	}

	cl.chats = append(cl.chats, &chat)
	cl.nextID += 1
	cl.open[owner] = cl.chats[chat.ID]
	return chat.ID
}

func (cl *ChatList) OpenChat(user string, id int) (*Chat, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if id < 0 || id >= len(cl.chats) {
		return nil, fmt.Errorf("invalid chat id")
	}
	cl.init()
	openChat := cl.chats[id]
	if openChat == nil {
		return nil, fmt.Errorf("chat[%d] does not exist", id)
	}
	if !openChat.isOwner(user) && !openChat.isUserInChat(user) {
		return nil, fmt.Errorf("user[%s] is not in chat[%d]", user, id)
	}
	cl.open[user] = cl.chats[id]
	return openChat, nil
}

func (cl *ChatList) CloseChat(user string, chatID int) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	userAtChat := cl.open[user]
	if userAtChat == nil {
		return fmt.Errorf("user[%s] has no open chat", user)
	}
	if userAtChat.ID != chatID {
		return fmt.Errorf("user[%s] is not open on chat[%d]", user, chatID)
	}
	cl.open[user] = nil
	return nil
}

func (cl *ChatList) GetOpenChat(user string) *Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	return cl.open[user]
}

func (cl *ChatList) GetChats(user string) []*Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	var userChats []*Chat
	for _, chat := range cl.chats {
		if chat == nil {
			continue
		}
		if chat.isOwner(user) || chat.isUserInChat(user) {
			userChats = append(userChats, chat)
		}
	}
	return userChats
}

func (cl *ChatList) GetChat(user string, chatID int) (*Chat, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	var chat *Chat
	for _, chat = range cl.chats {
		if chat == nil {
			continue
		}
		if chat.ID == chatID {
			break
		}
	}

	if chat == nil {
		return nil, fmt.Errorf("chatID[%d] does not exist", chatID)
	}
	if !chat.isOwner(user) && !chat.isUserInChat(user) {
		return nil, fmt.Errorf("user[%s] is not in chat[%d]", user, chatID)
	}

	return chat, nil
}

func (cl *ChatList) DeleteChat(user string, chat *Chat) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	if chat == nil {
		return fmt.Errorf("user[%s] cannot remove NIL chat", user)
	}
	if !chat.isOwner(user) {
		return fmt.Errorf("user[%s] is not owner of chat %d", user, chat.ID)
	}
	if cl.chats[chat.ID] == nil {
		return fmt.Errorf("chat[%d] is NIL", chat.ID)
	}
	cl.chats[chat.ID] = nil

	if cl.open[user] == nil {
		// noop
	} else if cl.open[user].ID == chat.ID {
		cl.open[user] = nil
	}

	return nil
}

func (cl *ChatList) InviteUser(user string, chatID int, invitee string) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	if chatID < 0 || chatID >= len(cl.chats) {
		return fmt.Errorf("invalid chat index[%d]", chatID)
	}
	chat := cl.chats[chatID]
	if !chat.isOwner(user) {
		return fmt.Errorf("user[%s] is not owner of chat %d", user, chatID)
	}
	err := chat.AddUser(user, invitee)
	if err != nil {
		return err
	}
	return nil
}
