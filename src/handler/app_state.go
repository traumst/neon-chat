package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"prplchat/src/db"
	"prplchat/src/model/app"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

var ApplicationState AppState

// TODO track users
type AppState struct {
	mu       sync.Mutex
	isInit   bool
	chats    app.ChatList
	userConn ActiveConnections
	config   utils.Config
}

func (state *AppState) Init(db *db.DBConn, config utils.Config) {
	ApplicationState = AppState{
		isInit:   true,
		chats:    app.ChatList{},
		userConn: make(ActiveConnections, 0),
		config:   config,
	}
}

func (state *AppState) LoadLocal() bool {
	return state.config.LoadLocal
}

func (state *AppState) SmtpConfig() utils.SmtpConfig {
	return state.config.Smtp
}

// CONN
func (state *AppState) AddConn(w http.ResponseWriter, r http.Request, user *app.User) *Conn {
	state.mu.Lock()
	defer state.mu.Unlock()

	if state.userConn == nil {
		log.Printf("[%s] AppState.AddConn TRACE init UserConn\n", h.GetReqId(&r))
		state.userConn = make(ActiveConnections, 0)
	}

	log.Printf("[%s] AppState.AddConn TRACE add conn for user[%d]\n", h.GetReqId(&r), user.Id)
	return state.userConn.Add(user, h.GetReqId(&r), w, r)
}

func (state *AppState) GetConn(userId uint) []*Conn {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetConn TRACE get conn for user[%d]\n", userId)
	return state.userConn.Get(userId)
}

func (state *AppState) DropConn(conn *Conn) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("[%s] AppState.DropConn TRACE drop conn[%s] user[%d]\n",
		conn.Origin, conn.Origin, conn.User.Id)
	return state.userConn.Drop(conn)
}

// USER
func (state *AppState) InviteUser(userId uint, chatId int, invitee *app.User) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.InviteUser TRACE invite user[%d] chat[%d] by user[%d]\n", invitee.Id, chatId, userId)
	if userId == invitee.Id {
		return fmt.Errorf("user cannot invite self")
	}
	return state.chats.InviteUser(userId, chatId, invitee)
}

func (state *AppState) DropUser(userId uint, chatId int, removeId uint) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.DropUser TRACE removing user[%d] chat[%d] by user[%d]\n", removeId, chatId, userId)
	_ = state.chats.CloseChat(removeId, chatId)
	return state.chats.ExpelUser(userId, chatId, removeId)
}

// CHAT
func (state *AppState) AddChat(user *app.User, chatName string) int {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.AddChat TRACE add chat[%s] for user[%d]\n", chatName, user.Id)
	chatId := state.chats.AddChat(user, chatName)
	return chatId
}

func (state *AppState) GetChats(userId uint) []*app.Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetChats TRACE get chats for user[%d]\n", userId)
	return state.chats.GetChats(userId)
}

func (state *AppState) GetChat(userId uint, chatId int) (*app.Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetChat TRACE get chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.GetChat(userId, chatId)
}

func (state *AppState) OpenChat(userId uint, chatId int) (*app.Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.OpenChat TRACE open chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.OpenChat(userId, chatId)
}

func (state *AppState) GetOpenChat(userId uint) *app.Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetOpenChat TRACE get open chat for user[%d]\n", userId)
	return state.chats.GetOpenChat(userId)
}

func (state *AppState) CloseChat(userId uint, chatId int) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.CloseChat TRACE close chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.CloseChat(userId, chatId)
}

func (state *AppState) DeleteChat(userId uint, chat *app.Chat) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.DeleteChat TRACE get chats for user[%d]\n", userId)
	_ = state.chats.CloseChat(userId, chat.Id)
	return state.chats.DeleteChat(userId, chat)
}
