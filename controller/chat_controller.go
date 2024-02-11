package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"go.chat/model"
	"go.chat/utils"
)

type ChatController struct {
	mu      sync.Mutex
	counter *atomic.Int32
	isInit  bool
}

func (c *ChatController) init() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isInit {
		return
	}
	c.isInit = true

	log.Printf("------ ChatController.init TRACE\n")
	c.counter = &atomic.Int32{}
	c.counter.Store(0)
}

func (c *ChatController) IsAnyoneConnected() bool {
	return c.counter.Load() > 0
}

func (c *ChatController) OpenChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> OpenChat\n", utils.GetReqId(r))
	if r.Method != "GET" {
		log.Printf("<-%s-- OpenChat TRACE auth does not allow %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR auth, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	path := utils.ParseUrlPath(r)
	log.Printf("--%s-> OpenChat, %s\n", utils.GetReqId(r), path[2])
	chatID, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR id, %s\n", utils.GetReqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}
	if chatID < 0 {
		log.Printf("<-%s-- OpenChat ERROR chatID, %s\n", utils.GetReqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", utils.GetReqId(r), chatID)
	openChat, err := app.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR chat, %s\n", utils.GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", utils.GetReqId(r))
	html, err := openChat.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR html template, %s\n", utils.GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- OpenChat TRACE returning template\n", utils.GetReqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func (c *ChatController) AddChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddChat\n", utils.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("<-%s-- AddChat TRACE auth does not allow %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Printf("--%s-> AddChat TRACE check login\n", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR auth, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	chatName := r.FormValue("chatName")
	log.Printf("--%s-> AddChat TRACE adding chat[%s]\n", utils.GetReqId(r), chatName)
	chatID := app.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE opening chat[%s][%d]\n", utils.GetReqId(r), chatName, chatID)
	openChat, err := app.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR chat, %s\n", utils.GetReqId(r), err)
		errMsg := fmt.Sprintf("ERROR: %s", err.Error())
		w.Write([]byte(errMsg))
		return
	}

	conn := app.getConn(user)
	if conn == nil {
		log.Printf("--%s-> AddChat ERROR cannot distribute chat header[%s] to user[%s]\n",
			utils.GetReqId(r), openChat.Log(), user)
	} else {
		log.Printf("--%s-> AddChat TRACE distributing chat header[%s] to user[%s]\n",
			utils.GetReqId(r), openChat.Log(), user)
		conn.Channel <- model.UserUpdate{
			Type: model.ChatUpdate,
			Chat: openChat,
			Msg:  nil,
			User: user,
		}
	}

	log.Printf("--%s-> AddChat TRACE templating chat[%s][%d]\n", utils.GetReqId(r), chatName, chatID)
	html, err := openChat.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("<--%s-- AddChat ERROR html, %s", utils.GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- AddChat TRACE writing response\n", utils.GetReqId(r))

	w.WriteHeader(http.StatusPartialContent)
	w.Write([]byte(html))
}

func (c *ChatController) InviteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> InviteUser\n", utils.GetReqId(r))
	if r.Method != "POST" {
		log.Printf("<-%s-- InviteUser TRACE auth does not allow %s\n", utils.GetReqId(r), r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> InviteUser TRACE check login\n", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR auth, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	chatID, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR chat id, %s\n", utils.GetReqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}
	invitee := r.FormValue("invitee")
	log.Printf("--%s-> InviteUser TRACE inviting[%s] to chat[%d]\n", utils.GetReqId(r), invitee, chatID)
	err = app.InviteUser(user, chatID, invitee)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR invite, %s\n", utils.GetReqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- InviteUser TRACE user %s added to chat [%d]\n", utils.GetReqId(r), invitee, chatID)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(fmt.Sprintf(" [%s] ", invitee)))
}
