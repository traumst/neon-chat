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

func (state *AppState) AddConn(w http.ResponseWriter, r http.Request, user string) *Conn {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("---%s--> AppState.AddConn TRACE add conn for user[%s]\n", utils.GetReqId(&r), user)
	if state.userConn == nil {
		state.userConn = make(UserConn, 0)
	}
	conn := state.userConn.Add(user, utils.GetReqId(&r), w, r)

	return conn
}

func (state *AppState) GetConn(reqId string, user string) *Conn {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞---%s---> AppState.GetConn TRACE get conn for user[%s]\n", reqId, user)
	conn, err := state.userConn.Get(reqId, user)
	if err != nil {
		return nil
	}

	return conn
}

func (state *AppState) DropConn(reqId string, conn *Conn, user string) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞---%s---> AppState.DropConn TRACE drop conn for user[%s] from [%s]\n", reqId, user, conn.Origin)
	state.userConn.Drop(reqId, conn)
}

func (state *AppState) AddChat(reqId string, user string, chatName string) int {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞---%s---> AppState.AddChat TRACE user[%s] chat[%s]\n", reqId, user, chatName)
	chatID := state.chats.AddChat(user, chatName)
	return chatID
}

func (state *AppState) OpenChat(user string, chatID int) (*Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.OpenChat TRACE user[%s] chat[%d]\n", user, chatID)
	return state.chats.OpenChat(user, chatID)
}

func (state *AppState) InviteUser(user string, chatID int, invitee string) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.InviteUser TRACE user[%s] chat[%d] invitee[%s]\n", user, chatID, invitee)
	return state.chats.InviteUser(user, chatID, invitee)
}

func (state *AppState) CloseChat(user string) {
	panic("implement closing chat")
}

func (state *AppState) GetChats(user string) []*Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.GetChats TRACE user[%s]\n", user)
	return state.chats.GetChats(user)
}

func (state *AppState) GetChat(user string, chatID int) *Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.OpenChat TRACE user[%s] chat[%d]\n", user, chatID)
	return state.chats.GetChat(user, chatID)
}

func (state *AppState) GetOpenChat(user string) *Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.GetOpenChat TRACE user[%s]\n", user)
	return state.chats.GetOpenChat(user)
}
