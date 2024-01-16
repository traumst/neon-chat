package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.chat/handlers"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Chat\n", reqId(r))
	_, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		OpenChat(w, r)
	case "POST":
		AddChat(w, r)
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func OpenChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> OpenChat\n", reqId(r))
	_, err := r.Cookie("username")
	if err != nil {
		log.Printf("--%s-> OpenChat ERROR username\n", reqId(r))
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/chat/"))
	if err != nil {
		log.Printf("--%s-> OpenChat ERROR id\n", reqId(r))
		handlers.SendBack(w, r)
		return
	}

	openChat := chats.OpenChat(chatID)
	html, err := openChat.GetHTML()
	if err != nil {
		log.Printf("--%s-> OpenChat ERROR html, %s\n", reqId(r), err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(html))
}

func AddChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> AddChat\n", reqId(r))
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	chatName := r.FormValue("chatName")
	chats.AddChat(usernameCookie.Value, chatName)

	http.Redirect(w, r, "/", http.StatusFound)
}
