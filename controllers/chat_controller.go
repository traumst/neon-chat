package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"go.chat/models"
	"go.chat/utils"
)

const stop = 60

type ChatController struct {
	//mu      sync.Mutex
	//clients map[*Client]bool
	counter *atomic.Int32
}

// type Client struct {
// 	send chan string
// }

func (c *ChatController) IsAnyoneConnected() bool {
	return c.counter.Load() > 0
}

func (c *ChatController) PollChats(w http.ResponseWriter, r *http.Request) {
	reqId := SetReqId(r) // TODO have to set manually without a middleware
	if c.counter == nil {
		c.counter = &atomic.Int32{}
		c.counter.Store(1)
	} else {
		c.counter.Add(1)
	}
	defer c.counter.Add(-1)
	log.Printf("--%s-> PollChats TRACE IN as %dth\n", reqId, c.counter.Load())

	if r.Method != "GET" {
		log.Printf("<-%s-- PollChats TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- PollChats ERROR auth, %s\n", reqId, err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	// Set the content type to text/event-stream
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// Create a new ticker that ticks every 5 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	// Count the number of loops
	loopCount := 0
	// Infinite loop to send SSE messages
	for {
		select {
		case <-ticker.C:
			log.Printf("<-%s-- PollChats TRACE loop SSE, %d client/s connected\n", GetReqId(r), c.counter.Load())
			loopCount += 1
			if loopCount >= stop {
				log.Printf("<-%s-- PollChats WARN stop SSE on %dth loop\n", GetReqId(r), loopCount)
				return
			}
			message := fmt.Sprintf("<li id=\"chat-%d\">chat-%d</li>", loopCount, loopCount)
			log.Printf("<-%s-- PollChats TRACE send SSE: [%s]\n", GetReqId(r), message)
			fmt.Fprintf(w, "id: %d\n\n", loopCount)
			fmt.Fprintf(w, "event: %s\n", models.ChatEventName)
			fmt.Fprintf(w, "data: %s\n\n", message)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			log.Printf("<-%s-- PollChats INFO conn closed\n", GetReqId(r))
			return
		}
	}
}

func (c *ChatController) OpenChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> OpenChat\n", GetReqId(r))
	if r.Method != "GET" {
		log.Printf("<-%s-- OpenChat TRACE auth does not allow %s\n", GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR auth, %s\n", GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	path := utils.ParseUrlPath(r)
	log.Printf("--%s-> OpenChat, %s\n", GetReqId(r), path[2])
	chatID, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR id, %s\n", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}
	if chatID < 0 {
		log.Printf("<-%s-- OpenChat ERROR chatID, %s\n", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", GetReqId(r), chatID)
	openChat, err := chats.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR chat, %s\n", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", GetReqId(r))
	html, err := openChat.GetHTML()
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR html template, %s\n", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- OpenChat TRACE returning template\n", GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func (c *ChatController) AddChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddChat\n", GetReqId(r))
	if r.Method != "POST" {
		log.Printf("<-%s-- AddChat TRACE auth does not allow %s\n", GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> AddChat TRACE check login\n", GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR auth, %s\n", GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	chatName := r.FormValue("chatName")
	log.Printf("--%s-> AddChat TRACE adding chat[%s]\n", GetReqId(r), chatName)
	chatID := chats.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE opening chat[%s][%d]\n", GetReqId(r), chatName, chatID)
	openChatTemplate, err := chats.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR chat, %s\n", GetReqId(r), err)
		errMsg := fmt.Sprintf("ERROR: %s", err.Error())
		w.Write([]byte(errMsg))
		return
	}
	html, err := openChatTemplate.GetHTML()
	if err != nil {
		log.Printf("<--%s-- AddChat ERROR html, %s", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- AddChat TRACE swriting response\n", GetReqId(r))

	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func (c *ChatController) InviteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> InviteUser\n", GetReqId(r))
	if r.Method != "POST" {
		log.Printf("<-%s-- InviteUser TRACE auth does not allow %s\n", GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> InviteUser TRACE check login\n", GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR auth, %s\n", GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	chatID, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR chat id, %s\n", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}
	invitee := r.FormValue("invitee")
	log.Printf("--%s-> InviteUser TRACE inviting[%s] to chat[%d]\n", GetReqId(r), invitee, chatID)
	inviteeUser, err := chats.InviteUser(user, chatID, invitee)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR invite, %s\n", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- InviteUser TRACE redirecting %s\n", GetReqId(r), inviteeUser)
	http.Redirect(w, r, fmt.Sprintf("/chat/%d", chatID), http.StatusFound)
}
