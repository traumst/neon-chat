package app

import (
	"fmt"
	"sync"
)

type ChatList struct {
	mu     sync.Mutex
	chats  []*Chat
	open   map[uint]*Chat
	nextId int
	isInit bool
}

func (cl *ChatList) init() {
	if cl.isInit {
		return
	}

	cl.isInit = true
	cl.chats = []*Chat{}
	cl.open = make(map[uint]*Chat)
}

func (cl *ChatList) AddChat(owner *User, chatName string) int {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	chat := Chat{
		Id:      cl.nextId,
		Name:    chatName,
		Owner:   owner,
		users:   []*User{owner},
		history: MessageStore{},
	}

	cl.chats = append(cl.chats, &chat)
	cl.nextId += 1
	cl.open[owner.Id] = cl.chats[chat.Id]
	return chat.Id
}

func (cl *ChatList) OpenChat(userId uint, chatId int) (*Chat, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if chatId < 0 || chatId >= len(cl.chats) {
		return nil, fmt.Errorf("invalid chat id")
	}
	cl.init()
	openChat := cl.chats[chatId]
	if openChat == nil {
		return nil, fmt.Errorf("chat[%d] does not exist", chatId)
	}
	if !openChat.isOwner(userId) && !openChat.isUserInChat(userId) {
		return nil, fmt.Errorf("user[%d] is not in chat[%d]", userId, chatId)
	}
	cl.open[userId] = cl.chats[chatId]
	return openChat, nil
}

func (cl *ChatList) CloseChat(userId uint, chatId int) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	userAtChat := cl.open[userId]
	if userAtChat == nil {
		return fmt.Errorf("user[%d] has no open chat", userId)
	}
	if userAtChat.Id != chatId {
		return fmt.Errorf("user[%d] is not open on chat[%d]", userId, chatId)
	}
	cl.open[userId] = nil
	return nil
}

func (cl *ChatList) GetOpenChat(userId uint) *Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	return cl.open[userId]
}

func (cl *ChatList) GetChats(userId uint) []*Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	var userChats []*Chat
	for _, chat := range cl.chats {
		if chat == nil {
			continue
		}
		if chat.isOwner(userId) || chat.isUserInChat(userId) {
			userChats = append(userChats, chat)
		}
	}
	return userChats
}

func (cl *ChatList) GetChat(userId uint, chatId int) (*Chat, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	var chat *Chat
	for _, chat = range cl.chats {
		if chat == nil {
			continue
		}
		if chat.Id == chatId {
			break
		}
	}

	if chat == nil {
		return nil, fmt.Errorf("chatId[%d] does not exist", chatId)
	}
	if !chat.isOwner(userId) && !chat.isUserInChat(userId) {
		return nil, fmt.Errorf("user[%d] is not in chat[%d]", userId, chatId)
	}

	return chat, nil
}

func (cl *ChatList) DeleteChat(userId uint, chat *Chat) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	if chat == nil {
		return fmt.Errorf("user[%d] cannot remove NIL chat", userId)
	}
	if !chat.isOwner(userId) {
		return fmt.Errorf("user[%d] is not owner of chat %d", userId, chat.Id)
	}
	if chat.Id >= len(cl.chats) || chat.Id < 0 {
		return fmt.Errorf("chat[%d] is out of range", chat.Id)
	}
	if cl.chats[chat.Id] == nil {
		return fmt.Errorf("chat[%d] is NIL", chat.Id)
	}
	cl.chats[chat.Id] = nil

	if cl.open[userId] == nil {
		// noop
	} else if cl.open[userId].Id == chat.Id {
		cl.open[userId] = nil
	}

	return nil
}

func (cl *ChatList) InviteUser(userId uint, chatId int, invitee *User) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	if chatId < 0 || chatId >= len(cl.chats) {
		return fmt.Errorf("invalid chat index[%d]", chatId)
	}
	chat := cl.chats[chatId]
	if !chat.isOwner(userId) {
		return fmt.Errorf("user[%d] is not owner of chat %d", userId, chatId)
	}
	err := chat.AddUser(userId, invitee)
	if err != nil {
		return err
	}
	return nil
}

func (cl *ChatList) ExpelUser(userId uint, chatId int, removeId uint) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.init()
	if chatId < 0 || chatId >= len(cl.chats) {
		return fmt.Errorf("invalid chat index[%d]", chatId)
	}
	chat := cl.chats[chatId]
	if chat == nil {
		return fmt.Errorf("chat[%d] not found", chatId)
	}
	if !chat.isOwner(userId) && userId != removeId {
		return fmt.Errorf("user[%d] is not owner of chat %d", userId, chatId)
	}
	err := chat.RemoveUser(userId, removeId)
	if err != nil {
		return err
	}
	return nil
}
