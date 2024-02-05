package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.chat/model"
	"go.chat/utils"
)

type UserConn map[string]Conn

type Conn struct {
	User   string
	Origin string
	Conn   http.ResponseWriter
}

func (uc *UserConn) IsConn(user string) bool {
	_, ok := (*uc)[user]
	return ok
}

func (uc *UserConn) Add(user string, origin string, conn http.ResponseWriter) {
	(*uc)[user] = Conn{
		User:   user,
		Origin: origin,
		Conn:   conn,
	}
}

func (uc *UserConn) Remove(user string) {
	delete(*uc, user)
}

var chats = model.ChatList{}
var userConns = make(UserConn)

var userUpdateChannel = make(chan *model.UserUpdate, 256)

const pingFreq = 5 * time.Second

var lastPing = time.Now()

func pollUpdates(w http.ResponseWriter, r *http.Request, user string) {
	if !userConns.IsConn(user) {
		log.Printf("<-%s-- pollUpdates WARN user not connected, %s\n", utils.GetReqId(r), user)
		userConns.Add(user, utils.GetReqId(r), w)
		defer userConns.Remove(user)
	}

	log.Printf("<-%s-- pollUpdates TRACE called by %s\n", utils.GetReqId(r), user)
	utils.SetSseHeaders(w)
	loopUpdates(w, r, user)
}

func loopUpdates(w http.ResponseWriter, r *http.Request, user string) {
	log.Printf("<-%s-- loopUpdates TRACE loopin, triggered by %s\n", utils.GetReqId(r), user)

	tick := time.NewTicker(pingFreq)
	for {
		select {
		case <-r.Context().Done():
			err := r.Context().Err()
			if err != nil {
				log.Printf("<-%s-- loopUpdates WARN conn closed, %s\n", utils.GetReqId(r), err)
			} else {
				log.Printf("<-%s-- loopUpdates INFO conn closed\n", utils.GetReqId(r))
			}
			return
		case <-tick.C:
			pingPong(w, r, user)
		case update := <-userUpdateChannel:
			updateUser(w, r, update, user)
		}

	}
}

func updateUser(
	w http.ResponseWriter,
	r *http.Request,
	up *model.UserUpdate,
	user string) {
	log.Printf("--%s-- updateUser TRACE IN, user[%s], input[%s]\n", utils.GetReqId(r), user, up.Log())
	if up.User == "" || up.Chat == nil {
		log.Printf("--%s-- updateUser INFO msg is empty, %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	switch up.Type {
	case model.ChatUpdate:
		sendChat(w, r, up, user)
	case model.MessageUpdate:
		sendMessage(w, r, up, user)
	default:
		log.Printf("--%s-- updateUser ERROR unknown update type, %s\n", utils.GetReqId(r), up.Log())
	}
}

func sendChat(
	w http.ResponseWriter,
	r *http.Request,
	up *model.UserUpdate,
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
	distribute(w, r, up.Chat, model.ChatEventName, user, eventID, html)
	log.Printf("--%s-- sendChat TRACE message from %s\n", utils.GetReqId(r), user)
}

func sendMessage(
	w http.ResponseWriter,
	r *http.Request,
	up *model.UserUpdate,
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
	w http.ResponseWriter,
	r *http.Request,
	chat *model.Chat,
	msgType model.SSEvent,
	eventID string,
	user string,
	html string) {
	log.Printf("---- distribute TRACE IN %s-%s by user %s\n", msgType, eventID, user)

	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")

	chatUsers, err := chat.GetUsers(user)
	if err != nil || chatUsers == nil {
		log.Printf("---- distribute TRACE get users, %s\n", err)
		return
	}

	for _, u := range chatUsers {
		if u == user && msgType == model.MessageEventName {
			continue
		}

		send(&w, r, msgType, user, eventID, html)
	}
	log.Printf("---- distribute TRACE OUT %s-%s by user %s\n", msgType, eventID, user)
	lastPing = time.Now()
}

func pingPong(w http.ResponseWriter, r *http.Request, user string) {
	if time.Since(lastPing) > pingFreq {
		log.Printf("--%s-> pingPong TRACE ping\n", utils.GetReqId(r))
		eventID := fmt.Sprintf("%s-%s", model.PingEventName, utils.RandStringBytes(5))
		data := fmt.Sprintf("data: %s\n\n", time.Now().Format(time.RFC3339))
		send(&w, r, model.PingEventName, user, eventID, data)
	}
	lastPing = time.Now()
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
