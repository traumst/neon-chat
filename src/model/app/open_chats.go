package app

import (
	"fmt"
	"log"
	"sync"

	"prplchat/src/utils/store"
)

type OpenChats struct {
	mu    sync.Mutex
	chats *store.LRUCache
	// userId -> chat
	open map[uint]*Chat
}

func NewOpenChats() *OpenChats {
	cl := OpenChats{
		chats: store.NewLRUCache(1024),
		open:  make(map[uint]*Chat),
	}
	return &cl
}

func (cl *OpenChats) GetChat(userId uint, chatId uint) (*Chat, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.getCached(chatId, userId)
}
func (cl *OpenChats) GetChats(userId uint) []*Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	chatIds := cl.chats.Keys()
	var userChats []*Chat
	for _, chatId := range chatIds {
		chat, _ := cl.getCached(chatId, userId)
		if chat == nil {
			continue
		}
		userChats = append(userChats, chat)
	}
	return userChats
}

func (cl *OpenChats) AddChat(chatId uint, owner *User, chatName string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	chat := Chat{
		Id:      chatId,
		Name:    chatName,
		Owner:   owner,
		users:   []*User{owner},
		history: MessageStore{},
	}
	cl.chats.Set(chatId, &chat)
	cl.open[owner.Id] = &chat
}

// fails if chat does not exist
func (cl *OpenChats) OpenChat(userId uint, chatId uint) (*Chat, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	openChat, err := cl.getCached(chatId, userId)
	if err != nil {
		return nil, err
	}
	cl.open[userId] = openChat
	return openChat, nil
}

func (cl *OpenChats) GetOpenChat(userId uint) *Chat {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.open[userId]
}

func (cl *OpenChats) CloseChat(userId uint, chatId uint) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
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

func (cl *OpenChats) DeleteChat(userId uint, chat *Chat) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if chat == nil {
		return fmt.Errorf("user[%d] cannot remove NIL chat", userId)
	}
	if !chat.isOwner(userId) {
		return fmt.Errorf("user[%d] is not owner of chat %d", userId, chat.Id)
	}
	count, err := cl.deleteCached(chat.Id)
	if err != nil {
		return err
	}
	log.Printf("Closed chat chat[%d] for [%d] users on chat deletion", chat.Id, count)
	return nil
}

func (cl *OpenChats) GetUser(userId uint) (*User, error) {
	// TODO rethink
	for _, chatId := range cl.chats.Keys() {
		chat, _ := cl.GetChat(chatId, userId)
		if chat == nil {
			continue
		}
		chatUsers, err := chat.GetUsers(userId)
		if err != nil {
			log.Printf("GetUser ERROR getting users from chat[%d] for user[%d]: %s", chatId, userId, err)
			continue
		}
		for _, user := range chatUsers {
			if user.Id == userId {
				return user, nil
			}
		}
	}
	return nil, nil
}

func (cl *OpenChats) InviteUser(userId uint, chatId uint, invitee *User) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	chat, err := cl.getCached(chatId, userId)
	if err != nil || chat == nil {
		return err
	}
	if !chat.isOwner(userId) {
		return fmt.Errorf("user[%d] is not owner of chat %d", userId, chatId)
	}
	err = chat.AddUser(userId, invitee)
	if err != nil {
		return err
	}
	return nil
}

func (cl *OpenChats) ExpelUser(userId uint, chatId uint, removeId uint) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	chat, err := cl.getCached(chatId, userId)
	if chat == nil || err != nil {
		return fmt.Errorf("chat[%d] not found, [%s]", chatId, err)
	}
	err = chat.RemoveUser(userId, removeId)
	if err != nil {
		return err
	}
	return nil
}

func (cl *OpenChats) SyncUser(user *User) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	chatIds := cl.chats.Keys()
	for _, chatId := range chatIds {
		chat, _ := cl.getCached(chatId, user.Id)
		if chat == nil {
			continue
		}
		err := chat.SyncUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cl *OpenChats) getCached(chatId uint, userId uint) (*Chat, error) {
	cachedObj, err := cl.chats.Get(chatId)
	if err != nil {
		return nil, err
	}
	if cachedObj == nil {
		return nil, fmt.Errorf("chat[%d] is not cached", chatId)
	}
	openChat, ok := cachedObj.(*Chat)
	if !ok {
		return nil, fmt.Errorf("chat[%d] is not a chat but [%t]", chatId, cachedObj)
	}
	if !openChat.isOwner(userId) && !openChat.IsUserInChat(userId) {
		return nil, fmt.Errorf("user[%d] is not in chat[%d]", userId, chatId)
	}
	return openChat, nil
}

// removes chat and closes it for all users
func (cl *OpenChats) deleteCached(chatId uint) (int, error) {
	deleted, err := cl.chats.Take(chatId)
	if err != nil {
		return 0, err
	}
	if deleted == nil {
		return 0, nil
	}
	chat, ok := deleted.(*Chat)
	if !ok {
		return 0, fmt.Errorf("chat[%d] is not a chat", chatId)
	}
	closed := 0
	for _, user := range chat.users {
		if user == nil {
			continue
		}
		if cl.open[user.Id] == nil {
			continue
		}
		if cl.open[user.Id].Id == chat.Id {
			cl.open[user.Id] = nil
			closed += 1
		}
	}
	return closed, nil
}
