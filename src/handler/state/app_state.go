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
		return err
	}
	if removedIdWrap == nil {
		return nil
	}
	removedId, ok := removedIdWrap.(uint)
	if removedId != 0 && removedId != chatId {
		state.chats.Set(userId, removedId)
	}
	if !ok {
		return fmt.Errorf("failed to remove chat[%d] for user[%d], checked[%d]", chatId, userId, removedId)
	}
	return nil
}
