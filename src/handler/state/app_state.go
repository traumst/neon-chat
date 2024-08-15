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

var GlobalAppState State

type State struct {
	mu     sync.Mutex
	isInit bool
	// users  store.LRUCache
	chats  store.LRUCache
	conns  OpenConnections
	config utils.Config
}

func (state *State) Init(config utils.Config) {
	if config.CacheSize <= 0 {
		config.CacheSize = 1024
	}
	GlobalAppState = State{
		isInit: true,
		chats:  *store.NewLRUCache(config.CacheSize),
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

func (state *State) OpenChat(userId uint, chatId uint) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.OpenChat TRACE open chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.Set(userId, chatId)
}

func (state *State) GetOpenChat(userId uint) uint {
	state.mu.Lock()
	defer state.mu.Unlock()
	log.Printf("AppState.GetOpenChat TRACE get open chat for user[%d]\n", userId)
	chatIdWrap, err := state.chats.Get(userId)
	if err != nil {
		return 0
	}
	chatId, ok := chatIdWrap.(uint)
	if !ok {
		return 0
	}
	return chatId
}

func (state *State) CloseChat(userId uint, chatId uint) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.CloseChat TRACE close chat[%d] for user[%d]\n", chatId, userId)
	removedIdWrap, err := state.chats.Take(userId)
	if err != nil {
		return err
	}
	if removedIdWrap == nil {
		return nil
	}
	removedId, ok := removedIdWrap.(uint)
	if removedId != 0 && removedId != chatId {
		state.chats.Set(userId, removedId)
	}
	if ok {
		return nil
	}
	return fmt.Errorf("failed to remove chat[%d] for user[%d], checked[%d]", chatId, userId, removedId)
}
