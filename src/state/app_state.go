package state

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/utils/config"
	h "neon-chat/src/utils/http"
	"neon-chat/src/utils/store"
)

var GlobalAppState State

type State struct {
	mu     sync.Mutex
	isInit bool
	chats  store.LRUCache
	conns  OpenConnections
	config config.Config
}

func (state *State) Init(config config.Config) {
	if config.CacheSize <= 0 || config.CacheSize > consts.MaxCacheSize {
		config.CacheSize = consts.MaxCacheSize
	}
	GlobalAppState = State{
		isInit: true,
		chats:  *store.NewLRUCache(config.CacheSize),
		conns:  make(OpenConnections, 0),
		config: config,
	}
}

func (state *State) SmtpConfig() (*config.SmtpConfig, error) {
	if !state.isInit {
		return nil, fmt.Errorf("state is not initialized")
	}
	return &state.config.Smtp, nil
}

// CONN
func (state *State) AddConn(w http.ResponseWriter, r http.Request, user *app.User, openChat *app.Chat) *Conn {
	if !state.isInit {
		log.Printf("ERROR AddConn state is not initialized")
		return nil
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.conns.Add(user, h.GetReqId(&r), w, r)
}

func (state *State) GetConn(userId uint) []*Conn {
	if !state.isInit {
		log.Printf("ERROR GetConn state is not initialized")
		return nil
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.conns.Get(userId)
}

func (state *State) DropConn(conn *Conn) error {
	if !state.isInit {
		log.Printf("ERROR DropConn state is not initialized")
		return fmt.Errorf("state is not initialized")
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.conns.Drop(conn)
}

// CHAT
func (state *State) OpenChat(userId uint, chatId uint) error {
	if !state.isInit {
		log.Printf("ERROR OpenChat state is not initialized")
		return fmt.Errorf("state is not initialized")
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.chats.Set(userId, chatId)
}

func (state *State) GetOpenChat(userId uint) uint {
	if !state.isInit {
		log.Printf("ERROR GetOpenChat state is not initialized")
		return 0
	}
	state.mu.Lock()
	defer state.mu.Unlock()
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
	if !state.isInit {
		return fmt.Errorf("state is not initialized")
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	removedIdWrap, err := state.chats.Take(userId)
	if err != nil {
		return fmt.Errorf("failed to close open chat, %s", err.Error())
	}
	if removedIdWrap == nil {
		return nil
	}
	removedId, ok := removedIdWrap.(uint)
	if !ok || removedId != 0 && removedId != chatId {
		state.chats.Set(userId, removedId)
		return nil
	}
	return nil
}
