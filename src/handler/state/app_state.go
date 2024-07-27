package state

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"prplchat/src/model/app"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
	"prplchat/src/utils/store"
)

var Application State

type State struct {
	mu     sync.Mutex
	isInit bool
	users  store.LRUCache
	chats  app.HotChats
	conns  OpenConnections
	config utils.Config
}

func (state *State) Init(config utils.Config) {
	cacheSize := config.CacheSize
	if cacheSize <= 0 {
		cacheSize = 1024
	}
	Application = State{
		isInit: true,
		users:  *store.NewLRUCache(128),
		chats:  *app.NewHotChats(),
		conns:  make(OpenConnections, 0),
		config: config,
	}
}

func (state *State) SmtpConfig() utils.SmtpConfig {
	return state.config.Smtp
}

// CONN
func (state *State) AddConn(w http.ResponseWriter, r http.Request, user *app.User, openChat *app.Chat) *Conn {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.conns == nil {
		panic("AddConn state.conns is nil")
	}

	return state.conns.Add(user, h.GetReqId(&r), w, r)
}

func (state *State) GetConn(userId uint) []*Conn {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.conns == nil {
		panic("GetConn state.conns is nil")
	}

	return state.conns.Get(userId)
}

func (state *State) DropConn(conn *Conn) error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.conns == nil {
		panic("DropConn state.conns is nil")
	}

	return state.conns.Drop(conn)
}

// USER
func (state *State) InviteUser(userId uint, chatId uint, invitee *app.User) error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if userId == invitee.Id {
		return fmt.Errorf("user cannot invite self")
	}

	return state.users.Set(invitee.Id, invitee)
}

func (state *State) GetUser(userId uint) (*app.User, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetUser TRACE retrieving user[%d]\n", userId)
	cached, err := state.users.Get(userId)
	if err != nil {
		return nil, fmt.Errorf("user[%d] not found: %s", userId, err.Error())
	}
	return cached.(*app.User), nil
}

func (state *State) ExpelFromChat(userId uint, chatId uint, removeId uint) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.ExpelFromChat TRACE user[%d] expels user[%d] from chat[%d]\n", userId, removeId, chatId)
	return state.chats.ExpelUser(userId, chatId, removeId)
}

func (state *State) UpdateUser(userId uint, template app.User) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.UpdateUser TRACE updating user[%d], template[%v]\n", userId, template)
	return state.chats.SyncUser(&template)
}

// CHAT
func (state *State) AddChat(chatId uint, chatName string, chatOwner *app.User) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.AddChat TRACE add chat[%d][%s] for user[%d]\n", chatId, chatName, chatOwner.Id)
	state.chats.AddChat(chatId, chatOwner, chatName)
	return nil
}

func (state *State) GetChats(userId uint) []*app.Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetChats TRACE get chats for user[%d]\n", userId)
	return state.chats.GetChats(userId)
}

func (state *State) GetChat(userId uint, chatId uint) (*app.Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetChat TRACE get chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.GetChat(userId, chatId)
}

func (state *State) OpenChat(userId uint, chatId uint) (*app.Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.OpenChat TRACE open chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.OpenChat(userId, chatId)
}

func (state *State) GetOpenChat(userId uint) *app.Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetOpenChat TRACE get open chat for user[%d]\n", userId)
	return state.chats.GetOpenChat(userId)
}

func (state *State) CloseChat(userId uint, chatId uint) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.CloseChat TRACE close chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.CloseChat(userId, chatId)
}

func (state *State) DeleteChat(userId uint, chat *app.Chat) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.DeleteChat TRACE get chats for user[%d]\n", userId)
	_ = state.chats.CloseChat(userId, chat.Id)
	return state.chats.DeleteChat(userId, chat)
}
