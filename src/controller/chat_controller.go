package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
)

func Welcome(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] Welcome TRACE\n", reqId)
	user, _ := handler.ReadSession(state, db, w, r)
	html, err := handler.TemplateWelcome(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to template html, %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func OpenChat(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] OpenChat TRACE\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] OpenChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		var msg string
		if err != nil {
			msg = err.Error()
		}
		log.Printf("[%s] OpenChat INFO user is not authorized, %s\n", h.GetReqId(r), msg)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	log.Printf("[%s] OpenChat, %s\n", reqId, path[2])
	chatId, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("[%s] OpenChat INFO invalid chat-id, %s\n", reqId, err)
		Welcome(state, db, w, r)
		return
	}
	if chatId < 0 {
		log.Printf("[%s] OpenChat ERROR chatId, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid chat id %d", chatId)))
		return
	}

	log.Printf("[%s] OpenChat TRACE chat[%d]\n", reqId, chatId)
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

func AddChat(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] AddChat TRACE\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] AddChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("action not allowed"))
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if user == nil {
		log.Printf("[%s] AddChat INFO user is not authorized, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user is not authorized"))
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

	log.Printf("[%s] AddChat TRACE adding user[%d] chat[%s]\n", reqId, user.Id, chatName)
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

func CloseChat(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] CloseChat\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] CloseChat WARN user, %s\n", h.GetReqId(r), err)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] CloseChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

func DeleteChat(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] DeleteChat TRACE\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] DeleteChat WARN user, %s\n", h.GetReqId(r), err)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] DeleteChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
