package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.chat/utils"
)

func OpenChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> OpenChat\n", reqId(r))
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("--%s-> OpenChat ERROR auth, %s\n", reqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	path := utils.ParseUrlPath(r)
	log.Printf("--%s-> OpenChat, %s\n", reqId(r), path[2])
	chatID, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("--%s-> OpenChat ERROR id, %s\n", reqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", reqId(r), chatID)
	openChat, err := chats.OpenChat(user, chatID)
	if err != nil {
		log.Printf("--%s-> OpenChat ERROR chat, %s\n", reqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", reqId(r))
	html, err := openChat.GetHTML()
	if err != nil {
		log.Printf("--%s-> OpenChat ERROR html template, %s\n", reqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}

	log.Printf("<-%s-- OpenChat TRACE returning template\n", reqId(r))
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func AddChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddChat\n", reqId(r))
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> AddChat TRACE check login\n", reqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR auth, %s\n", reqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	chatName := r.FormValue("chatName")
	log.Printf("--%s-> AddChat TRACE adding chat[%s]\n", reqId(r), chatName)
	chatID := chats.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE opening chat[%s][%d]\n", reqId(r), chatName, chatID)
	openChatTemplate, err := chats.OpenChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR chat, %s\n", reqId(r), err)
		errMsg := fmt.Sprintf("ERROR: %s", err.Error())
		w.Write([]byte(errMsg))
		return
	}
	html, err := openChatTemplate.GetHTML()
	if err != nil {
		log.Printf("<--%s-- AddChat ERROR html, %s", reqId(r), err)
		utils.SendBack(w, r, http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- AddChat TRACE swriting response\n", reqId(r))

	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}
