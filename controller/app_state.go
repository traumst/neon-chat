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

var chats = model.ChatList{}
var clients = make([]*model.Client, 0)
var userUpdateChannel = make(chan *model.UserUpdate, 1024)

var pingFreq = 5 * time.Second

func pollUpdates(w http.ResponseWriter, r *http.Request, user string) {
	registerUser(w, r, user)

	utils.SetSseHeaders(w)

	tick := time.NewTicker(pingFreq)
	lastMsg := time.Now()
	pingPong(w, r, lastMsg)
	for {
		select {
		case <-r.Context().Done():
			err := r.Context().Err()
			if err != nil {
				log.Printf("<-%s-- pollUpdates WARN conn closed, %s\n", utils.GetReqId(r), err)
			} else {
				log.Printf("<-%s-- pollUpdates INFO conn closed\n", utils.GetReqId(r))
			}
			return
		case <-tick.C:
			pingPong(w, r, lastMsg)
		case update := <-userUpdateChannel:
			updateUser(w, r, update, user)
		}

		lastMsg = time.Now()
	}
}

func registerUser(w http.ResponseWriter, r *http.Request, user string) {
	var known *model.Client
	for _, c := range clients {
		if c.User == user {
			known = c
			break
		}
	}

	if known != nil {
		log.Printf("--%s-> updateClientList TRACE client already known\n", utils.GetReqId(r))
		known.Request = r
		known.Response = &w
	} else {
		log.Printf("--%s-> updateClientList TRACE adding client\n", utils.GetReqId(r))
		clients = append(clients, &model.Client{User: user, Request: r, Response: &w})
	}
}

func pingPong(w http.ResponseWriter, r *http.Request, lastPing time.Time) {
	if time.Since(lastPing) > pingFreq {
		log.Printf("--%s-> pingPong TRACE ping\n", utils.GetReqId(r))
		fmt.Fprintf(w, "id: %s-%s\n\n", model.PingEventName, utils.RandStringBytes(5))
		fmt.Fprintf(w, "event: %s\n", model.PingEventName)
		fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))

		w.(http.Flusher).Flush()
	}
}

func updateUser(
	w http.ResponseWriter,
	r *http.Request,
	up *model.UserUpdate,
	user string) {
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

	log.Printf("--%s-- sendChat TRACE templating %s\n",
		utils.GetReqId(r), up.Msg.Log())
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
	w http.ResponseWriter,
	r *http.Request,
	up *model.UserUpdate,
	user string) {
	if up.Msg == nil {
		log.Printf("--%s-- sendMessage INFO msg is empty, %s\n", utils.GetReqId(r), up.Msg.Log())
		return
	}

	log.Printf("--%s-- sendMessage TRACE templating %s\n",
		utils.GetReqId(r), up.Msg.Log())
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
	log.Printf("--%s-- distribute TRACE event %s-%s by user %s\n", utils.GetReqId(r), msgType, eventID, user)

	// must escape newlines in SSE
	html = strings.ReplaceAll(html, "\n", " ")

	chatUsers, err := chat.GetUsers(user)
	if err != nil {
		log.Printf("--%s-- distribute ERROR get users, %s\n", utils.GetReqId(r), err)
		return
	}

	for _, u := range chatUsers {
		if u == user && msgType == model.MessageEventName {
			continue
		}

		// TODO get proper client and send them

		log.Printf("--%s-- distribute TRACE serving html to fucking %s\n", utils.GetReqId(r), u)
		fmt.Fprintf(w, "id: %s\n\n", eventID)
		fmt.Fprintf(w, "event: %s\n", model.MessageEventName)
		fmt.Fprintf(w, "data: %s\n\n", html)
	}
	w.(http.Flusher).Flush()
}
