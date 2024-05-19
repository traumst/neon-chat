package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"go.chat/src/db"
	"go.chat/src/model/app"
	"go.chat/src/utils"
	h "go.chat/src/utils/http"
)

var ApplicationState AppState

type AppState struct {
	mu       sync.Mutex
	isInit   bool
	db       *db.DBConn
	chats    app.ChatList
	userConn UserConn
	users    []app.User
	config   AppConfig
}

type AppConfig struct {
	LoadLocal bool
	Smtp      SmtpConfig
}

type SmtpConfig struct {
	User string
	Pass string
	Host string
	Port string
}

func (state *AppState) Init(db *db.DBConn, config AppConfig) {
	ApplicationState = AppState{
		isInit:   true,
		db:       db,
		chats:    app.ChatList{},
		userConn: make(UserConn, 0),
		users:    make([]app.User, 0),
		config:   config,
	}
}

func (state *AppState) LoadLocal() bool {
	return state.config.LoadLocal
}

func (state *AppState) SmtpConfig() SmtpConfig {
	return state.config.Smtp
}

// CONN
func (state *AppState) ReplaceConn(w http.ResponseWriter, r http.Request, user *app.User) *Conn {
	conn, err := state.GetConn(user.Id)
	for err == nil && conn != nil {
		log.Printf("[%s] AppState.ReplaceConn WARN drop old conn to user[%d]\n", h.GetReqId(&r), user.Id)
		state.DropConn(conn)
		conn, err = state.GetConn(user.Id)
	}

	return state.addConn(w, r, user)
}

func (state *AppState) addConn(w http.ResponseWriter, r http.Request, user *app.User) *Conn {
	state.mu.Lock()
	defer state.mu.Unlock()

	if state.userConn == nil {
		log.Printf("[%s] AppState.AddConn TRACE init UserConn\n", h.GetReqId(&r))
		state.userConn = make(UserConn, 0)
	}

	log.Printf("[%s] AppState.AddConn TRACE add conn for user[%d]\n", h.GetReqId(&r), user.Id)
	return state.userConn.Add(user, h.GetReqId(&r), w, r)
}

func (state *AppState) GetConn(userId uint) (*Conn, error) {
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
// TODO untrack
func (state *AppState) TrackUser(user *app.User) error {
	if user == nil {
		return fmt.Errorf("user was nil")
	}
	for _, u := range state.users {
		if u.Id == user.Id {
			log.Printf("AppState.TrackUser TRACE tracked user[%d] update", user.Id)
			u.Name = user.Name
			for _, c := range state.chats.GetChats(user.Id) {
				c.SyncUser(user)
			}
			return nil
		}
	}
	log.Printf("AppState.TrackUser TRACE will track user[%d]\n", user.Id)
	state.users = append(state.users, *user)
	return nil
}

func (state *AppState) GetUser(userId uint) (*app.User, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetUser TRACE user[%d]\n", userId)
	user, err := state.db.GetUser(userId)
	if err != nil {
		return nil, err
	}
	log.Printf("AppState.GetUser TRACE user[%d] found[%s]\n", userId, user.Name)
	appUser := UserFromDB(*user)
	state.TrackUser(&appUser)
	return &appUser, nil
}

func (state *AppState) UpdateUser(appUser *app.User) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetUser UpdateUser user[%d]\n", appUser.Id)
	dbUser := UserToDB(*appUser)
	err := state.db.UpdateName(dbUser)
	if err != nil {
		return fmt.Errorf("failed to update user [%d], %s", appUser.Id, err.Error())
	}
	return state.TrackUser(appUser)
}

func (state *AppState) InviteUser(userId uint, chatId int, invitee *app.User) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.InviteUser TRACE invite user[%d] chat[%d] by user[%d]\n",
		invitee.Id, chatId, userId)
	if userId == invitee.Id {
		return fmt.Errorf("user cannot invite self")
	}
	return state.chats.InviteUser(userId, chatId, invitee)
}

func (state *AppState) DropUser(userId uint, chatId int, removeId uint) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.DropUser TRACE removing user[%d] chat[%d] by user[%d]\n", removeId, chatId, userId)
	// remove from tracked
	for i, u := range state.users {
		if u.Id == removeId {
			state.users = append(state.users[:i], state.users[i+1:]...)
			break
		}
	}
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

func (state *AppState) GetOpenChat(userId uint) *app.Chat {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetOpenChat TRACE get open chat for user[%d]\n", userId)
	return state.chats.GetOpenChat(userId)
}

func (state *AppState) OpenChat(userId uint, chatId int) (*app.Chat, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.OpenChat TRACE open chat[%d] for user[%d]\n", chatId, userId)
	return state.chats.OpenChat(userId, chatId)
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

func (state *AppState) AddAvatar(userId uint, avatar *app.Avatar) (*app.Avatar, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.AddAvatar TRACE user[%d], avatar[%s]\n", userId, avatar.Title)
	dbAvatar, err := state.db.AddAvatar(userId, avatar.Title, avatar.Image, avatar.Mime)
	if err != nil {
		return nil, fmt.Errorf("avatar not added: %s", err)
	}
	return &app.Avatar{
		Id:     dbAvatar.Id,
		UserId: dbAvatar.UserId,
		Title:  dbAvatar.Title,
		Size:   fmt.Sprintf("%dKB", dbAvatar.Size/utils.KB),
		Image:  dbAvatar.Image,
		Mime:   dbAvatar.Mime,
	}, nil
}

func (state *AppState) GetAvatar(userId uint) (*app.Avatar, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetAvatar TRACE user[%d]\n", userId)
	dbAvatar, err := state.db.GetAvatar(userId)
	if err != nil {
		return nil, fmt.Errorf("avatar not found: %s", err)
	}
	return &app.Avatar{
		Id:     dbAvatar.Id,
		UserId: dbAvatar.UserId,
		Title:  dbAvatar.Title,
		Size:   fmt.Sprintf("%dKB", dbAvatar.Size/utils.KB),
		Image:  dbAvatar.Image,
		Mime:   dbAvatar.Mime,
	}, nil
}

func (state *AppState) GetAvatars(userId uint) ([]*app.Avatar, error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.GetAvatar TRACE user[%d]\n", userId)
	dbAvatars, err := state.db.GetAvatars(userId)
	if err != nil {
		return nil, fmt.Errorf("avatar not found: %s", err)
	}
	avatars := make([]*app.Avatar, len(dbAvatars))
	for _, dbAvatar := range dbAvatars {
		avatars = append(avatars, &app.Avatar{
			Id:     dbAvatar.Id,
			UserId: dbAvatar.UserId,
			Title:  dbAvatar.Title,
			Size:   fmt.Sprintf("%dKB", dbAvatar.Size/utils.KB),
			Image:  dbAvatar.Image,
			Mime:   dbAvatar.Mime,
		})
	}
	return avatars, nil
}

func (state *AppState) DropAvatar(avatar *app.Avatar) error {
	if avatar == nil {
		return fmt.Errorf("cannot drop nil avatar")
	}
	state.mu.Lock()
	defer state.mu.Unlock()

	log.Printf("AppState.DropAvatar TRACE drops avatar[%d]\n", avatar.Id)
	return state.db.DropAvatar(avatar.Id)
}
