package model

import (
	"log"
	"net/http"
	"sync"

	"go.chat/utils"
)

type AppState struct {
	mu       sync.Mutex
	chats    ChatList
	userConn UserConn
}

func (state *AppState) ReplaceConn(w http.ResponseWriter, r http.Request, user string) *Conn {
	conn, err := state.GetConn(user)
	if err == nil && conn != nil {
		log.Printf("∞---%s---> AppState.ReplaceConn TRACE drop old conn to user[%s]\n", utils.GetReqId(&r), user)
		state.DropConn(conn, user)
	}

	return state.AddConn(w, r, user)
}

func (state *AppState) AddConn(w http.ResponseWriter, r http.Request, user string) *Conn {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.userConn == nil {
		log.Printf("----%s---> AppState.AddConn TRACE init UserConn\n", utils.GetReqId(&r))
		state.userConn = make(UserConn, 0)
	}

	log.Printf("----%s---> AppState.AddConn TRACE add conn for user[%s]\n", utils.GetReqId(&r), user)
	return state.userConn.Add(user, utils.GetReqId(&r), w, r)
}

func (state *AppState) GetConn(user string) (*Conn, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	return state.userConn.Get(user)
}

func (state *AppState) DropConn(conn *Conn, user string) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	return state.userConn.Drop(conn)
}

func (state *AppState) AddChat(user string, chatName string) int {
	state.mu.Lock()
	defer state.mu.Unlock()

	chatID := state.chats.AddChat(user, chatName)
	return chatID
}

func (state *AppState) OpenChat(user string, chatID int) (*Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	return state.chats.OpenChat(user, chatID)
}

func (state *AppState) InviteUser(user string, chatID int, invitee string) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	return state.chats.InviteUser(user, chatID, invitee)
}

func (state *AppState) CloseChat(reqId string, user string) error {
	defer recover()
	panic("implement closing chat")
}

func (state *AppState) GetChats(user string) []*Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	return state.chats.GetChats(user)
}

func (state *AppState) GetChat(user string, chatID int) (*Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	return state.chats.GetChat(user, chatID)
}

func (state *AppState) GetOpenChat(user string) *Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	return state.chats.GetOpenChat(user)
}
