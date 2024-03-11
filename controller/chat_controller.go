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

func OpenChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> OpenChat\n", reqId)
	if r.Method != "GET" {
		log.Printf("<-%s-- OpenChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR auth, %s\n", reqId, err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	path := utils.ParseUrlPath(r)
	log.Printf("--%s-> OpenChat, %s\n", reqId, path[2])
	chatID, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if chatID < 0 {
		log.Printf("<-%s-- OpenChat ERROR chatID, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", reqId, chatID)
	openChat, err := app.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR chat, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", reqId)
	html, err := openChat.ToTemplate(user).GetHTML()
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR html template, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- OpenChat TRACE returning template\n", reqId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func AddChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> AddChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- AddChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Printf("--%s-> AddChat TRACE check login\n", reqId)
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR auth, %s\n", reqId, err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	chatName := r.FormValue("chatName")
	log.Printf("--%s-> AddChat TRACE adding user[%s] chat[%s]\n", reqId, user, chatName)
	chatID := app.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE user[%s] opening chat[%s][%d]\n", reqId, user, chatName, chatID)
	openChat, err := app.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR chat, %s\n", reqId, err)
		errMsg := fmt.Sprintf("ERROR: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}

	log.Printf("--%s-> AddChat TRACE templating chat[%s][%d]\n",
		reqId, chatName, chatID)

	sendChatContent(reqId, w, openChat, user)
	err = handler.DistributeChat(app, openChat, user, user, model.ChatCreated)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR cannot distribute chat header, %s\n",
			reqId, err)
	}
}

func InviteUser(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> InviteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- InviteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Printf("--%s-> InviteUser TRACE check login\n", reqId)
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR auth, %s\n", reqId, err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	chatID, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invitee := r.FormValue("invitee")
	log.Printf("--%s-> InviteUser TRACE inviting[%s] to chat[%d]\n", reqId, invitee, chatID)
	err = app.InviteUser(user, chatID, invitee)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	chat, err := app.GetChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR cannot find chat[%d], %s\n", reqId, chatID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.DistributeChat(app, chat, user, invitee, model.ChatInvite)

	log.Printf("<-%s-- InviteUser TRACE user [%s] added to chat [%d] by user [%s]\n",
		reqId, invitee, chatID, user)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(fmt.Sprintf(" [%s] ", invitee)))
}

func CloseChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> CloseChat\n", reqId)
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatID := r.PostFormValue("chatid")
	if chatID == "" {
		log.Printf("<-%s-- CloseChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chatID)
	if err != nil {
		log.Printf("<-%s-- CloseChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.CloseChat(user, id)
	if err != nil {
		log.Printf("<-%s-- CloseChat ERROR close chat[%d] for [%s], %s\n",
			reqId, id, user, err)
	}
	welcome := model.WelcomeTemplate{ActiveUser: user}
	html, err := welcome.GetHTML()
	if err != nil {
		log.Printf("<-%s-- CloseChat ERROR cannot template welcome page, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- CloseChat TRACE user[%s] closed chat [%d]\n",
		reqId, user, id)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteChat TRACE\n", reqId)
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatID := r.PostFormValue("chatid")
	if chatID == "" {
		log.Printf("<-%s-- DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chatID)
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user, id)
	if err != nil || chat == nil {
		log.Printf("<-%s-- DeleteChat ERROR cannot get chat[%d] for [%s]\n", reqId, id, user)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = app.DeleteChat(user, chat)
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR remove chat[%d] from [%s], %s\n", reqId, id, chat.Name, err)
	}

	handler.DistributeChat(app, chat, user, user, model.ChatClose)
	handler.DistributeChat(app, chat, user, user, model.ChatDeleted)

	log.Printf("<-%s-- DeleteChat TRACE user[%s] deletes chat [%d]\n", reqId, user, id)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte("[DELETED]"))
}

func sendChatContent(reqId string, w http.ResponseWriter, chat *model.Chat, user string) {
	template := chat.ToTemplate(user)
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
