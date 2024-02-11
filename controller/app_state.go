package controller

import (
	"log"
	"net/http"
	"sync"

	"go.chat/handler"
	"go.chat/model"
	"go.chat/utils"
)

var app = App{}

type AppState struct {
	chats    model.ChatList
	userConn model.UserConn
}

type App struct {
	mu    sync.Mutex
	state AppState
}

func (app *App) AddChat(user string, chatName string) int {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.state.chats.AddChat(user, chatName)
}

func (app *App) OpenChat(user string, chatID int) (*model.Chat, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.state.chats.OpenChat(user, chatID)
}

func (app *App) InviteUser(user string, chatID int, invitee string) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.state.chats.InviteUser(user, chatID, invitee)
}

func (app *App) GetChats(user string) []*model.Chat {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.state.chats.GetChats(user)
}

func (app *App) GetOpenChat(user string) *model.Chat {
	app.mu.Lock()
	defer app.mu.Unlock()
	return app.state.chats.GetOpenChat(user)
}

func (app *App) PollUpdates(w http.ResponseWriter, r http.Request, user string) {
	log.Printf("--%s-> PollUpdates TRACE polling updates for [%s]\n", utils.GetReqId(&r), user)
	conn := app.storeConn(w, r, user)
	if conn == nil {
		log.Printf("<-%s-- PollUpdates ERROR conn not be established for [%s]\n", utils.GetReqId(&r), user)
		return
	}

	log.Printf("--%s-> PollUpdates TRACE sse initiated for [%s]\n", utils.GetReqId(&r), user)
	utils.SetSseHeaders(w)
	app.consumeUpdates(conn, user)
}

func (app *App) getConn(user string) *model.Conn {
	app.mu.Lock()
	defer app.mu.Unlock()
	conn, err := app.state.userConn.Get(user)
	if err != nil {
		return nil
	}
	return conn
}

func (app *App) storeConn(w http.ResponseWriter, r http.Request, user string) *model.Conn {
	log.Printf("∞-%s-∞ storeConn TRACE add conn for user[%s]\n", utils.GetReqId(&r), user)
	app.mu.Lock()
	defer app.mu.Unlock()
	if app.state.userConn == nil {
		app.state.userConn = make(model.UserConn, 0)
	}
	conn := app.state.userConn.Add(user, utils.GetReqId(&r), w, r)

	return conn
}

func (app *App) consumeUpdates(conn *model.Conn, user string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.loopUpdates(conn, user)
	}()
	wg.Wait()
}

func (app *App) loopUpdates(conn *model.Conn, user string) {
	log.Printf("∞-%s-∞ loopUpdates TRACE IN, triggered by [%s]\n", utils.GetReqId(&conn.Reader), conn.User)
	for range make([]int, 256) {
		select {
		case <-conn.Reader.Context().Done():
			err := conn.Reader.Context().Err()
			if err != nil {
				log.Printf("<-%s-∞ loopUpdates WARN conn closed from [%s], %s\n", utils.GetReqId(&conn.Reader), conn.User, err)
			} else {
				log.Printf("<-%s-∞ loopUpdates INFO conn closed [%s]\n", utils.GetReqId(&conn.Reader), conn.User)
			}
			return
		case update := <-conn.Channel:
			log.Printf("<-%s-∞ loopUpdates INFO update received for [%s] update [%s]\n",
				utils.GetReqId(&conn.Reader), conn.User, update.Log())
			app.sendUpdates(update, user, utils.GetReqId(&conn.Reader))
		}
	}
	log.Printf("---- loopUpdates TRACE OUT, triggered by [%s]\n", conn.User)
}

func (app *App) sendUpdates(up model.UserUpdate, user string, reqId string) {
	log.Printf("--%s-- sendUpdates TRACE IN, user[%s], input[%s]\n", reqId, user, up.Log())
	if up.User == "" || up.Chat == nil {
		log.Printf("---- updateUser INFO msg is empty, %s\n", up.Msg.Log())
		return
	}
	participants, err := up.Chat.GetUsers(user)
	if err != nil {
		log.Printf("--%s-- sendUpdates ERROR failed to get user[%s] chats, %s\n", reqId, user, err)
		return
	} else if participants == nil {
		log.Printf("--%s-- sendUpdates ERROR user[%s] has no chats\n", reqId, user)
		return
	}

	for _, p := range participants {
		if p == user && up.Type == model.MessageUpdate {
			log.Printf("--%s-- sendUpdates INFO skip sending message to origin sender [%s], update[%s]\n",
				reqId, user, up.Log())
			continue
		}
		// TODO check if user has matching chat open

		go app.trySend(up, p)
	}
	log.Printf("--%s-- sendUpdates TRACE OUT, user[%s], input[%s]\n", reqId, user, up.Log())
}

func (app *App) trySend(up model.UserUpdate, participant string) {
	log.Printf("---- trySend TRACE IN\n")
	conn := app.getConn(participant)
	if conn == nil {
		log.Printf("---- trySend ERROR user[%s] has no connection\n", participant)
		return
	}
	go app.informParticipant(conn, up, participant)
	log.Printf("---- trySend TRACE OUT\n")
}

func (app *App) informParticipant(conn *model.Conn, up model.UserUpdate, user string) {
	r := conn.Reader
	w := conn.Writer

	switch up.Type {
	case model.ChatUpdate:
		handler.SendChat(&w, up, user)
	case model.MessageUpdate:
		handler.SendMessage(&w, &r, up, user)
	default:
		log.Printf("--%s-- informParticipant ERROR unknown update type, %s\n", utils.GetReqId(&r), up.Log())
	}
}
