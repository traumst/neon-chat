package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/shared"
	"neon-chat/src/handler/state"
	a "neon-chat/src/model/app"
	"neon-chat/src/utils"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(utils.ReqIdKey).(string)
	log.Printf("[%s] Welcome TRACE\n", reqId)
	user := r.Context().Value(utils.ActiveUser).(*a.User)
	html, err := handler.TemplateWelcome(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to template html, %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func OpenChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(utils.ReqIdKey).(string)
	log.Printf("[%s] OpenChat TRACE\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] OpenChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}

	path := strings.Split(r.URL.Path, "/")
	log.Printf("[%s] OpenChat, %s\n", reqId, path[2])
	chatId, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("[%s] OpenChat INFO invalid chat-id, %s\n", reqId, err)
		Welcome(w, r)
		return
	}
	if chatId < 0 {
		log.Printf("[%s] OpenChat ERROR chatId, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid chat id %d", chatId)))
		return
	}

	user := r.Context().Value(utils.ActiveUser).(*a.User)
	state := r.Context().Value(utils.AppState).(*state.State)
	db := r.Context().Value(utils.DBConn).(*d.DBConn)
	html, err := handler.HandleChatOpen(state, db, user, uint(chatId))
	if err != nil {
		log.Printf("[%s] OpenChat ERROR cannot open chat[%d], %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to open chat"))
		return
	}

	log.Printf("[%s] OpenChat TRACE returning template\n", reqId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func AddChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(utils.ReqIdKey).(string)
	log.Printf("[%s] AddChat TRACE\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] AddChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("action not allowed"))
		return
	}

	chatName := r.FormValue("chatName")
	chatName = utils.ReplaceWithSingleSpace(chatName)
	chatName = utils.RemoveSpecialChars(chatName)
	if len(chatName) < 4 {
		log.Printf("[%s] AddChat ERROR chat name [%s]\n", reqId, chatName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("bad chat name [%s]", chatName)))
		return
	}

	user := r.Context().Value(utils.ActiveUser).(*a.User)
	state := r.Context().Value(utils.AppState).(*state.State)
	db := r.Context().Value(utils.DBConn).(*d.DBConn)
	html, err := handler.HandleChatAdd(state, db, user, chatName)
	if err != nil {
		log.Printf("sendChat ERROR cannot template chat[%s], %s", chatName, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to template chat"))
		return
	}

	log.Printf("[%s] AddChat TRACE chat[%s] created by user[%d]\n", reqId, chatName, user.Id)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func CloseChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(utils.ReqIdKey).(string)
	log.Printf("[%s] CloseChat\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatId, err := shared.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] CloseChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := r.Context().Value(utils.ActiveUser).(*a.User)
	state := r.Context().Value(utils.AppState).(*state.State)
	db := r.Context().Value(utils.DBConn).(*d.DBConn)
	html, err := handler.HandleChatClose(state, db, user, uint(chatId))
	if err != nil {
		log.Printf("[%s] CloseChat ERROR cannot template welcome page, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] CloseChat TRACE user[%d] closed chat[%d]\n", reqId, user.Id, chatId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(utils.ReqIdKey).(string)
	log.Printf("[%s] DeleteChat TRACE\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatId, err := shared.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] DeleteChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(utils.ActiveUser).(*a.User)
	state := r.Context().Value(utils.AppState).(*state.State)
	db := r.Context().Value(utils.DBConn).(*d.DBConn)
	err = handler.HandleChatDelete(state, db, user, chatId)
	if err != nil {
		log.Printf("[%s] DeleteChat WARN user[%d] failed to delete chat[%d], %s\n", reqId, user.Id, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] DeleteChat TRACE user[%d] deleted chat [%d]\n", reqId, user.Id, chatId)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("[DELETED_C]"))
}
