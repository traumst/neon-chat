package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.chat/handler"
	"go.chat/model"
	"go.chat/utils"
)

type ChatController struct{}

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if chatID < 0 {
		log.Printf("<-%s-- OpenChat ERROR chatID, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", utils.GetReqId(r), chatID)
	openChat, err := app.State.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR chat, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", utils.GetReqId(r))
	html, err := openChat.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR html template, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
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
	log.Printf("--%s-> AddChat TRACE adding user[%s] chat[%s]\n", utils.GetReqId(r), user, chatName)
	chatID := app.State.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE user[%s] opening chat[%s][%d]\n", utils.GetReqId(r), user, chatName, chatID)
	openChat, err := app.State.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR chat, %s\n", utils.GetReqId(r), err)
		errMsg := fmt.Sprintf("ERROR: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}

	log.Printf("--%s-> AddChat TRACE templating chat[%s][%d]\n",
		utils.GetReqId(r), chatName, chatID)
	template := openChat.ToTemplate(user)
	sendChatContent(utils.GetReqId(r), w, template)
	err = handler.DistributeChat(utils.GetReqId(r), &app.State, user, template, model.ChatCreated)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR cannot distribute chat header, %s\n",
			utils.GetReqId(r), err)
	}
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invitee := r.FormValue("invitee")
	log.Printf("--%s-> InviteUser TRACE inviting[%s] to chat[%d]\n", utils.GetReqId(r), invitee, chatID)
	err = app.State.InviteUser(user, chatID, invitee)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR invite, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chat, err := app.State.GetChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR cannot find chat[%d], %s\n", utils.GetReqId(r), chatID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temlate := chat.ToTemplate(invitee)
	handler.DistributeChat(utils.GetReqId(r), &app.State, invitee, temlate, model.ChatInvite)

	log.Printf("<-%s-- InviteUser TRACE user [%s] added to chat [%d] by user [%s]\n",
		utils.GetReqId(r), invitee, chatID, user)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(fmt.Sprintf(" [%s] ", invitee)))
}

func (c *ChatController) DeleteChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> DeleteChat\n", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR chat id, %s\n", utils.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.State.GetChat(user, id)
	if err != nil || chat == nil {
		log.Printf("<-%s-- DeleteMessage ERROR cannot get chat[%d] for [%s]\n", utils.GetReqId(r), id, user)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = app.State.DeleteChat(user, chat)
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR remove chat[%d] from [%s], %s\n",
			utils.GetReqId(r), id, chat.Name, err)
		// TODO not necessarily StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.DistributeChat(utils.GetReqId(r), &app.State, user, chat.ToTemplate(user), model.ChatDeleted)

	log.Printf("<-%s-- DeleteChat TRACE user[%s] PRETENDS to delete chat [%d]\n",
		utils.GetReqId(r), user, id)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(fmt.Sprintf("OK user[%s] chat[%d]", user, id)))
}

func sendChatContent(reqId string, w http.ResponseWriter, template *model.ChatTemplate) {
	html, err := template.GetHTML()
	if err != nil {
		log.Printf("<--%s-- sendChatContent ERROR cannot template chat [%+v], %s", reqId, template, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- sendChatContent TRACE writing response\n", reqId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}
