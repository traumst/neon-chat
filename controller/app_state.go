package controller

import (
	"log"
	"net/http"
	"sync"

	"go.chat/handler"
	"go.chat/model"
	"go.chat/utils"
)

var chats = model.ChatList{}
var userConns = make(model.UserConn)

func pollUpdates(w http.ResponseWriter, r http.Request, user string) {
	isConn, conn := userConns.IsConn(user)
	if !isConn {
		log.Printf("<-%s-- pollUpdates WARN user not connected, %s\n", utils.GetReqId(&r), user)
		conn = userConns.Add(user, utils.GetReqId(&r), w, r)
		defer userConns.Drop(user)
	}

	log.Printf("--%s-- pollUpdates TRACE called by %s\n", utils.GetReqId(&r), user)
	utils.SetSseHeaders(w)
	consumeUpdates(conn)
}

func consumeUpdates(conn *model.Conn) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loopUpdates(conn)
	}()
	wg.Wait()
}

func loopUpdates(conn *model.Conn) {
	log.Printf("---- loopUpdates TRACE IN, triggered by %s\n", conn.User)

	for {
		select {
		case <-conn.Reader.Context().Done():
			err := conn.Reader.Context().Err()
			if err != nil {
				log.Printf("<--- loopUpdates WARN conn closed, %s, %s\n", conn.User, err)
			} else {
				log.Printf("<--- loopUpdates INFO conn closed %s\n", conn.User)
			}
			return
		case update := <-conn.Channel:
			updateUser(&conn.Writer, &conn.Reader, update, conn.User)
		}
	}
}

func updateUser(
	w *http.ResponseWriter,
	r *http.Request,
	up model.UserUpdate,
	user string) {
	log.Printf("--%s-- updateUser TRACE IN, user[%s], input[%s]\n", utils.GetReqId(r), user, up.Log())
	if up.User == "" || up.Chat == nil {
		log.Printf("--%s-- updateUser INFO msg is empty, %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	switch up.Type {
	case model.ChatUpdate:
		handler.SendChat(w, r, up, user)
	case model.MessageUpdate:
		handler.SendMessage(w, r, up, user)
	default:
		log.Printf("--%s-- updateUser ERROR unknown update type, %s\n", utils.GetReqId(r), up.Log())
	}
}
