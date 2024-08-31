package state

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"neon-chat/src/model/app"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
	"neon-chat/src/utils/store"
)

var GlobalAppState State

type State struct {
	mu     sync.Mutex
	isInit bool
	chats  store.LRUCache
	conns  OpenConnections
	config utils.Config
}

func (state *State) Init(config utils.Config) {
	if config.CacheSize <= 0 || config.CacheSize > utils.MaxCacheSize {
		config.CacheSize = utils.MaxCacheSize
	}
	GlobalAppState = State{
		isInit: true,
		chats:  *store.NewLRUCache(config.CacheSize),
		conns:  make(OpenConnections, 0),
		config: config,
	}
}

func (state *State) SmtpConfig() (*utils.SmtpConfig, error) {
	if !state.isInit {
		return nil, fmt.Errorf("state is not initialized")
	}
	return &state.config.Smtp, nil
}

// CONN
func (state *State) AddConn(w http.ResponseWriter, r http.Request, user *app.User, openChat *app.Chat) *Conn {
	if !state.isInit {
		log.Printf("AddConn ERROR state is not initialized")
		return nil
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.conns.Add(user, h.GetReqId(&r), w, r)
}

func (state *State) GetConn(userId uint) []*Conn {
	if !state.isInit {
		log.Printf("GetConn ERROR state is not initialized")
		return nil
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.conns.Get(userId)
}

func (state *State) DropConn(conn *Conn) error {
	if !state.isInit {
		log.Printf("DropConn ERROR state is not initialized")
		return fmt.Errorf("state is not initialized")
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.conns.Drop(conn)
}

func (state *State) OpenChat(userId uint, chatId uint) error {
	if !state.isInit {
		log.Printf("OpenChat ERROR state is not initialized")
		return fmt.Errorf("state is not initialized")
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.chats.Set(userId, chatId)
}

func (state *State) GetOpenChat(userId uint) uint {
	if !state.isInit {
		log.Printf("GetOpenChat ERROR state is not initialized")
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
		log.Printf("GetOpenChat ERROR state is not initialized")
		return fmt.Errorf("state is not initialized")
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	removedIdWrap, err := state.chats.Take(userId)
	if err != nil {
		log.Printf("CloseChat ERROR failed to close open chat for user[%d], %s\n", userId, err.Error())
		return fmt.Errorf("failed to close open chat, %s", err.Error())
	}
	if removedIdWrap == nil {
		return nil
	}
	removedId, ok := removedIdWrap.(uint)
	if !ok || removedId != 0 && removedId != chatId {
		log.Printf("CloseChat WARN failed to close chat[%d] for user[%d], checked[%d]\n", chatId, userId, removedId)
		state.chats.Set(userId, removedId)
		return fmt.Errorf("failed to remove open chat for user[%d]", userId)
	}
	return nil
}
