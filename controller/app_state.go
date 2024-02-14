package controller

import (
	"log"
	"net/http"

	"go.chat/handler"
	"go.chat/model"
	"go.chat/utils"
)

var app = App{}

type App struct {
	state model.AppState
}

func (app *App) PollUpdatesForUser(w http.ResponseWriter, r http.Request, user string) {
	log.Printf("---%s--> PollUpdatesForUser TRACE IN polling updates for [%s]\n", utils.GetReqId(&r), user)

	conn := app.state.AddConn(w, r, user)
	if conn == nil {
		log.Printf("<--%s--- PollUpdatesForUser ERROR conn not be established for [%s]\n", utils.GetReqId(&r), user)
		return
	}
	defer app.state.DropConn(utils.GetReqId(&r), conn, user)

	log.Printf("---%s--> PollUpdatesForUser TRACE sse initiated for [%s]\n", utils.GetReqId(&r), user)
	utils.SetSseHeaders(&conn.Writer)
	app.pollUpdatesForUser(conn, user)
	log.Printf("<--%s--- PollUpdatesForUser TRACE OUT polling updates for [%s]\n", utils.GetReqId(&r), user)
}

func (app *App) pollUpdatesForUser(conn *model.Conn, user string) {
	log.Printf("∞--%s--> loopUpdates TRACE IN, triggered by [%s]\n", utils.GetReqId(&conn.Reader), conn.User)
	//ping := time.NewTicker(5 * time.Second)
	for i := range make([]int, 256) {
		log.Printf("∞--%s--> loopUpdates TRACE loop [%d], triggered by [%s]\n",
			utils.GetReqId(&conn.Reader), i, conn.User)

		select {
		case <-conn.Reader.Context().Done():
			err := conn.Reader.Context().Err()
			log.Printf("<--%s--∞ loopUpdates WARN conn closed from [%s], %s\n",
				utils.GetReqId(&conn.Reader), conn.User, err)
			return

		case update := <-conn.Channel:
			app.sendUpdates(utils.GetReqId(&conn.Reader), update, user)

			// case <-ping.C:
			// 	wg.Add(1)
			// 	go func() {
			// 		defer wg.Done()
			// 		app.sendUpdates(
			// 			utils.GetReqId(&conn.Reader),
			// 			model.UserUpdate{
			// 				User: user,
			// 				Msg:  "ping",
			// 				Type: model.PingUpdate,
			// 			},
			// 			user)
			// 	}()
			// 	wg.Wait()
		}
	}
	log.Printf("<--%s--∞ loopUpdates TRACE OUT, triggered by [%s]\n", utils.GetReqId(&conn.Reader), conn.User)
}

func (app *App) sendUpdates(reqId string, up model.UserUpdate, user string) {
	log.Printf("∞--%s--> sendUpdates TRACE IN, user[%s], input[%s]\n", reqId, user, up)
	if up.User == "" || up.Msg == "" {
		log.Printf("<--%s--∞ sendUpdates INFO user or msg is empty, %s\n", reqId, up)
		return
	}

	if up.User == user && up.Type == model.MessageUpdate {
		log.Printf("<--%s--∞ sendUpdates INFO skip sending message to origin sender [%s], update[%s]\n",
			reqId, user, up)
		return
	}

	app.trySend(reqId, up, user)
	log.Printf("<--%s--∞ sendUpdates TRACE OUT, user[%s], update[%s]\n", reqId, user, up)
}

func (app *App) trySend(reqId string, up model.UserUpdate, user string) {
	log.Printf("∞--%s--> trySend TRACE IN user[%s]\n", reqId, user)
	conn := app.state.GetConn(reqId, user)
	if conn == nil {
		log.Printf("<--%s--∞ trySend ERROR user[%s] has no connection\n", reqId, user)
		return
	}

	r := conn.Reader
	w := conn.Writer

	switch up.Type {
	case model.ChatUpdate:
	case model.ChatInvite:
	case model.MessageUpdate:
	case model.PingUpdate:
		log.Printf("<--%s--∞ trySend TRACE sending event user[%s], w[%T]\n", reqId, user, w)
		handler.SendEvent(reqId, &w, up, user)
	default:
		log.Printf("<--%s--∞ trySend ERROR unknown update type[%s], %s\n", utils.GetReqId(&r), up.Type, up)
	}
	log.Printf("<--%s--∞ trySend TRACE OUT user[%s]\n", reqId, user)
}
