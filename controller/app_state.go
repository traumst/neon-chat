package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"go.chat/model"
	"go.chat/utils"
)

var chats = model.ChatList{}
var userConns = make(model.UserConn)

//var userUpdateChannel = make(chan *model.UserUpdate, 256)

func pollUpdates(w http.ResponseWriter, r http.Request, user string) {
	isConn, conn := userConns.IsConn(user)
	if !isConn {
		log.Printf("<-%s-- pollUpdates WARN user not connected, %s\n", utils.GetReqId(&r), user)
		conn = userConns.Add(user, utils.GetReqId(&r), w, r)
		defer userConns.Remove(user)
	}

	log.Printf("--%s-- pollUpdates TRACE called by %s\n", utils.GetReqId(&r), user)
	utils.SetSseHeaders(w)

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

	// disconnect after 100 loops to refresh connection
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

	if up.Type == model.ChatUpdate {
		sendChat(w, r, up, user)
	} else if up.Type == model.MessageUpdate {
		sendMessage(w, r, up, user)
	} else {
		log.Printf("--%s-- updateUser ERROR unknown update type, %s\n", utils.GetReqId(r), up.Log())
	}
}

func sendChat(
	w *http.ResponseWriter,
	r *http.Request,
	up model.UserUpdate,
	user string) {
	if up.Msg != nil {
		log.Printf("--%s-- sendChat INFO unexpected input msg %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	log.Printf("--%s-- sendChat TRACE templating %s\n", utils.GetReqId(r), up.Chat.Log())
	html, err := up.Chat.ToTemplate(user).GetShortHTML()
	if err != nil {
		log.Printf("--%s-- sendChat ERROR templating %s, %s\n",
			utils.GetReqId(r), up.Msg.Log(), err)
		return
	}

	eventID := fmt.Sprintf("%s-%s", model.MessageEventName, utils.RandStringBytes(5))
	distribute(w, r, up.Chat, model.ChatEventName, eventID, user, html)
	log.Printf("--%s-- sendChat TRACE message from %s\n", utils.GetReqId(r), user)
}

func sendMessage(
	w *http.ResponseWriter,
	r *http.Request,
	up model.UserUpdate,
	user string) {
	if up.Msg == nil {
		log.Printf("--%s-- sendMessage INFO msg is empty, %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	log.Printf("--%s-- sendMessage TRACE templating %s\n", utils.GetReqId(r), up.Msg.Log())
	html, err := up.Msg.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("--%s-- sendMessage ERROR templating %s, %s\n",
			utils.GetReqId(r), up.Msg.Log(), err)
		return
	}

	eventID := fmt.Sprintf("%s-%s", model.MessageEventName, utils.RandStringBytes(5))
	distribute(w, r, up.Chat, model.MessageEventName, eventID, user, html)
	log.Printf("--%s-- sendMessage TRACE message from %s\n", utils.GetReqId(r), user)
}

func distribute(
	w *http.ResponseWriter,
	r *http.Request,
	chat *model.Chat,
	msgType model.SSEvent,
	eventID string,
	user string,
	html string) {
	log.Printf("---- distribute TRACE IN type[%s], event[%s] by user [%s]\n", msgType, eventID, user)

	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")
	html = strings.ReplaceAll(html, "  ", " ")

	chatUsers, err := chat.GetUsers(user)
	if err != nil || chatUsers == nil {
		log.Printf("---- distribute ERROR get users, %s\n", err)
		return
	}

	for _, u := range chatUsers {
		if u == user && msgType == model.MessageEventName {
			continue
		}

		send(w, r, msgType, user, eventID, html)
	}
	log.Printf("---- distribute TRACE OUT msg[%s], event[%s] by user [%s]\n", msgType, eventID, user)
}

func send(
	w *http.ResponseWriter,
	r *http.Request,
	eventName model.SSEvent,
	user string,
	eventID string,
	html string) {
	log.Printf("--%s-- send TRACE html to %s\n", utils.GetReqId(r), user)
	fmt.Fprintf(*w, "id: %s\n\n", eventID)
	fmt.Fprintf(*w, "event: %s\n", eventName)
	fmt.Fprintf(*w, "data: %s\n\n", html)
	(*w).(http.Flusher).Flush()
}
