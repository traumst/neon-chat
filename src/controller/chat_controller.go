package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	d "prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func Welcome(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] OpenChat TRACE returning template\n", reqId)
	user, err := handler.ReadSession(app, db, w, r)
	var welcome template.WelcomeTemplate
	if user == nil {
		log.Printf("[%s] OpenChat TRACE user is not authorized, %s\n", h.GetReqId(r), err)
		welcome = template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
	} else {
		log.Printf("[%s] OpenChat TRACE user is authorized, %s\n", h.GetReqId(r), err)
		welcome = template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
	}
	html, err := welcome.HTML()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to template html, %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func OpenChat(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] OpenChat\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] OpenChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if user == nil {
		log.Printf("[%s] OpenChat INFO user is not authorized, %s\n", h.GetReqId(r), err)
		// &template.InfoMessage{
		// 	Header: "User is not authenticated",
		// 	Body:   "Your session has probably expired",
		// 	Footer: "Reload the page and try again",
		// }
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	path := strings.Split(r.URL.Path, "/")
	log.Printf("[%s] OpenChat, %s\n", reqId, path[2])
	chatId, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("[%s] OpenChat INFO invalid chat-id, %s\n", reqId, err)
		Welcome(app, db, w, r)
		return
	}
	if chatId < 0 {
		log.Printf("[%s] OpenChat ERROR chatId, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid chat id %d", chatId)))
		return
	}
	log.Printf("[%s] OpenChat TRACE chat[%d]\n", reqId, chatId)
	var html string
	openChat, err := app.OpenChat(user.Id, uint(chatId))
	if err != nil {
		log.Printf("[%s] OpenChat ERROR chat, %s\n", reqId, err)
		welcome := template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
		html, err = welcome.HTML()
	} else {
		log.Printf("[%s] OpenChat TRACE html template\n", reqId)
		html, err = openChat.Template(user, user).HTML()
	}
	if err != nil {
		log.Printf("[%s] OpenChat ERROR html template, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template chat"))
		return
	}
	log.Printf("[%s] OpenChat TRACE returning template\n", reqId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func AddChat(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] AddChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] AddChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("action not allowed"))
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
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
	openChat, err := handler.HandleChatAdd(app, db, user, chatName)
	if err != nil {
		log.Printf("[%s] AddChat ERROR chat, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to open new chat"))
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = sse.DistributeChat(app, openChat, user, user, user, event.ChatAdd)
		if err != nil {
			log.Printf("[%s] AddChat ERROR cannot distribute chat header, %s\n", reqId, err)
		}
	}()
	go func() {
		defer wg.Done()
		template := openChat.Template(user, user)
		html, err := template.HTML()
		if err != nil {
			log.Printf("[%s] sendChatContent ERROR cannot template chat[%d], %s", reqId, template.ChatId, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to template chat"))
			return
		}
		log.Printf("[%s] sendChatContent TRACE writing response\n", reqId)
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(html))
	}()
	wg.Wait()
}

func CloseChat(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] CloseChat\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] CloseChat WARN user, %s\n", h.GetReqId(r), err)
		// &template.InfoMessage{
		// 	Header: "User is not authenticated",
		// 	Body:   "Your session has probably expired",
		// 	Footer: "Reload the page and try again",
		// }
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	chatIdStr := r.PostFormValue("chatid")
	if chatIdStr == "" {
		log.Printf("[%s] CloseChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		log.Printf("[%s] CloseChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.CloseChat(user.Id, uint(chatId))
	if err != nil {
		log.Printf("[%s] CloseChat ERROR close chat[%d] for user[%d], %s\n",
			reqId, chatId, user.Id, err)
	}
	welcome := template.WelcomeTemplate{User: *user.Template(0, 0, 0)}
	html, err := welcome.HTML()
	if err != nil {
		log.Printf("[%s] CloseChat ERROR cannot template welcome page, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[%s] CloseChat TRACE user[%d] closed chat[%d]\n", reqId, user.Id, chatId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteChat(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] DeleteChat TRACE\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] DeleteChat WARN user, %s\n", h.GetReqId(r), err)
		// &template.InfoMessage{
		// 	Header: "User is not authenticated",
		// 	Body:   "Your session has probably expired",
		// 	Footer: "Reload the page and try again",
		// }
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	inChatId := r.PostFormValue("chatid")
	if inChatId == "" {
		log.Printf("[%s] DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currChatid, err := strconv.Atoi(inChatId)
	if err != nil {
		log.Printf("[%s] DeleteChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId := uint(currChatid)
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil || chat == nil {
		log.Printf("[%s] DeleteChat ERROR cannot get chat[%d] for user[%d]\n", reqId, chatId, user.Id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.HandleChatDelete(app, db, user.Id, chat)
	if err != nil {
		log.Printf("[%s] DeleteChat ERROR remove chat[%d] from [%s], %s\n", reqId, chatId, chat.Name, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Printf("[%s] DeleteChat TRACE user[%d] deleted chat [%d]\n", reqId, user.Id, chatId)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		distrbuteChatDelete(app, chat, user)
	}()
	go func(chatId uint, userId uint) {
		defer wg.Done()
		log.Printf("[%s] DeleteChat TRACE user[%d] deletes chat[%d]\n", reqId, userId, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[DELETED_C]"))
	}(chat.Id, user.Id)
	wg.Wait()
}

func distrbuteChatDelete(app *state.State, chat *a.Chat, user *a.User) {
	log.Printf("distrbuteChatDelete TRACE distributes user[%d] closes chat[%d]\n", user.Id, chat.Id)
	err := sse.DistributeChat(app, chat, user, nil, user, event.ChatClose)
	if err != nil {
		log.Printf("distrbuteChatDelete ERROR cannot distribute chat close, %s\n", err)
		return
	}
	log.Printf("distrbuteChatDelete TRACE distributes user[%d] deletes chat[%d]\n", user.Id, chat.Id)
	err = sse.DistributeChat(app, chat, user, nil, user, event.ChatDrop)
	if err != nil {
		log.Printf("distrbuteChatDelete ERROR cannot distribute chat deleted, %s\n", err)
		return
	}
	log.Printf("distrbuteChatDelete TRACE distributes user[%d] expel all from chat[%d]\n", user.Id, chat.Id)
	err = sse.DistributeChat(app, chat, user, nil, nil, event.ChatExpel)
	if err != nil {
		log.Printf("distrbuteChatDelete ERROR cannot distribute chat user expel, %s\n", err)
		return
	}
}
