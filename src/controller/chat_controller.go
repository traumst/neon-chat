package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"neon-chat/src/consts"
	"neon-chat/src/controller/shared"
	"neon-chat/src/db"
	"neon-chat/src/handler/parse"
	"neon-chat/src/handler/pub"
	"neon-chat/src/model/app"
	"neon-chat/src/model/event"
	"neon-chat/src/sse"
	"neon-chat/src/state"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] Welcome TRACE\n", reqId)
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	html, err := shared.TemplateWelcome(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to template html, %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func OpenChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
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

	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	openChat, err := pub.OpenChat(state, dbConn, user, uint(chatId))
	if err != nil {
		log.Printf("[%s] OpenChat ERROR cannot open chat[%d], %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to open chat"))
		return
	}
	appChatUsers, err := pub.GetChatUsers(dbConn.Conn, openChat.Id)
	if err != nil {
		log.Printf("[%s] OpenChat ERROR cannot get chat[%d] users, %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to load participants"))
		return
	}
	appChatMsgs, err := pub.GetChatMessages(dbConn.Conn, openChat.Id)
	if err != nil {
		log.Printf("[%s] OpenChat ERROR cannot get chat[%d] messages, %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to load messages"))
		return
	}
	tmpl := openChat.Template(user, user, appChatUsers, appChatMsgs)
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("[%s] OpenChat ERROR cannot template chat[%d], %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to template open chat"))
		return
	}
	log.Printf("[%s] OpenChat TRACE returning template\n", reqId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func AddChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
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

	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, err := pub.AddChat(state, dbConn, user, chatName)
	if err != nil {
		log.Printf("AddChat ERROR failed to add chat[%s], %s", chatName, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to template chat"))
		return
	}
	err = sse.DistributeChat(state, dbConn.Tx, appChat, user, user, user, event.ChatAdd)
	if err != nil {
		log.Printf("AddChat ERROR cannot distribute chat[%d] creation to user[%d]: %s",
			appChat.Id, user.Id, err.Error())
	}
	tmpl := appChat.Template(user, user, []*app.User{user}, []*app.Message{})
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("AddChat ERROR cannot template chat[%s], %s", chatName, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to template chat"))
		return
	}
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("[%s] AddChat TRACE chat[%s] created by user[%d]\n", reqId, chatName, user.Id)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func CloseChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] CloseChat\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] CloseChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	err = pub.CloseChat(state, dbConn, user, uint(chatId))
	if err != nil {
		log.Printf("[%s] CloseChat ERROR failed to close chat, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	html, err := shared.TemplateWelcome(user)
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
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] DeleteChat TRACE\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] DeleteChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	deletedChat, err := pub.DeleteChat(state, dbConn, user, chatId)
	if err != nil {
		log.Printf("[%s] DeleteChat WARN user[%d] failed to delete chat[%d], %s\n", reqId, user.Id, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = sse.DistributeChat(state, dbConn.Tx, deletedChat, user, nil, user, event.ChatClose)
	if err != nil {
		log.Printf("DeleteChat ERROR cannot distribute chat close, %s", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, deletedChat, user, nil, user, event.ChatDrop)
	if err != nil {
		log.Printf("DeleteChat ERROR cannot distribute chat deleted, %s", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, deletedChat, user, nil, nil, event.ChatExpel)
	if err != nil {
		log.Printf("DeleteChat ERROR cannot distribute chat user expel, %s", err.Error())
	}
	w.(*h.StatefulWriter).IndicateChanges()
	log.Printf("[%s] DeleteChat TRACE user[%d] deleted chat [%d]\n", reqId, user.Id, chatId)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("[DELETED_C]"))
}
