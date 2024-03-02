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
	for err == nil && conn != nil {
		log.Printf("∞---%s---> AppState.ReplaceConn WARN drop old conn to user[%s]\n", utils.GetReqId(&r), user)
		state.DropConn(conn, user)
		conn, err = state.GetConn(user)
	}

	return state.addConn(w, r, user)
}

func (state *AppState) addConn(w http.ResponseWriter, r http.Request, user string) *Conn {
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

	log.Printf("∞--------> AppState.GetConn TRACE get conn for user[%s]\n", user)
	return state.userConn.Get(user)
}

func (state *AppState) DropConn(conn *Conn, user string) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("----%s---> AppState.DropConn TRACE drop conn[%s][%s] for user[%s]\n",
		conn.Origin, conn.Origin, conn.User, user)
	return state.userConn.Drop(conn)
}

func (state *AppState) AddChat(user string, chatName string) int {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.AddChat TRACE add chat[%s] for user[%s]\n", user, chatName)
	chatID := state.chats.AddChat(user, chatName)
	return chatID
}

func (state *AppState) CloseChat(user string, chatID int) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.CloseChat TRACE close chat[%d] for user[%s]\n", chatID, user)
	return state.chats.CloseChat(user, chatID)
}

func (state *AppState) DeleteChat(user string, chat *Chat) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.DeleteChat TRACE get chats for user[%s]\n", user)
	return state.chats.DeleteChat(user, chat)
}

func (state *AppState) OpenChat(user string, chatID int) (*Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.OpenChat TRACE open chat[%d] for user[%s]\n", chatID, user)
	return state.chats.OpenChat(user, chatID)
}

func (state *AppState) InviteUser(user string, chatID int, invitee string) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.InviteUser TRACE invite user[%s] chat[%d] for user[%s]\n", invitee, chatID, user)
	return state.chats.InviteUser(user, chatID, invitee)
}

func (state *AppState) GetChats(user string) []*Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.GetChats TRACE get chats for user[%s]\n", user)
	return state.chats.GetChats(user)
}

func (state *AppState) GetChat(user string, chatID int) (*Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.GetChat TRACE get chat[%d] for user[%s]\n", chatID, user)
	return state.chats.GetChat(user, chatID)
}

func (state *AppState) GetOpenChat(user string) *Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("∞--------> AppState.GetOpenChat TRACE get open chat for user[%s]\n", user)
	return state.chats.GetOpenChat(user)
}
