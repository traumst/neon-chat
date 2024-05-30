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
	"prplchat/src/model/event"
	"prplchat/src/model/template"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func Welcome(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] OpenChat TRACE returning template\n", reqId)
	user, err := handler.ReadSession(app, db, w, r)
	var welcome template.WelcomeTemplate
	if user == nil {
		log.Printf("[%s] OpenChat TRACE user is not authorized, %s\n", h.GetReqId(r), err)
		welcome = template.WelcomeTemplate{User: *user.Template(-1, 0, 0)}
	} else {
		log.Printf("[%s] OpenChat TRACE user is authorized, %s\n", h.GetReqId(r), err)
		welcome = template.WelcomeTemplate{User: *user.Template(-1, 0, 0)}
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

func OpenChat(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/", http.StatusFound)
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
	openChat, err := app.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("[%s] OpenChat ERROR chat, %s\n", reqId, err)
		welcome := template.WelcomeTemplate{User: *user.Template(-1, 0, 0)}
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

func AddChat(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
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
	chatName = utils.TrimSpaces(chatName)
	chatName = utils.TrimSpecial(chatName)
	if chatName == "" || len(chatName) < 4 {
		log.Printf("[%s] AddChat ERROR chat name [%s]\n", reqId, chatName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("bad chat name [%s]", chatName)))
		return
	}
	log.Printf("[%s] AddChat TRACE adding user[%d] chat[%s]\n", reqId, user.Id, chatName)
	chatId := app.AddChat(user, chatName)
	log.Printf("[%s] AddChat TRACE user[%d] opening chat[%s][%d]\n", reqId, user.Id, chatName, chatId)
	openChat, err := app.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("[%s] AddChat ERROR chat, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to open new chat"))
		return
	}
	log.Printf("[%s] AddChat TRACE templating chat[%s][%d]\n", reqId, chatName, chatId)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = handler.DistributeChat(app, openChat, user, user, user, event.ChatAdd)
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

func CloseChat(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
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

		http.Redirect(w, r, "/", http.StatusFound)
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
	err = app.CloseChat(user.Id, chatId)
	if err != nil {
		log.Printf("[%s] CloseChat ERROR close chat[%d] for user[%d], %s\n",
			reqId, chatId, user.Id, err)
	}
	welcome := template.WelcomeTemplate{User: *user.Template(-1, 0, 0)}
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

func DeleteChat(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
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

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	chatId := r.PostFormValue("chatid")
	if chatId == "" {
		log.Printf("[%s] DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chatId)
	if err != nil {
		log.Printf("[%s] DeleteChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, id)
	if err != nil || chat == nil {
		log.Printf("[%s] DeleteChat ERROR cannot get chat[%d] for user[%d]\n", reqId, id, user.Id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = app.DeleteChat(user.Id, chat)
	if err != nil {
		log.Printf("[%s] DeleteChat ERROR remove chat[%d] from [%s], %s\n", reqId, id, chat.Name, err)
	} else {
		log.Printf("[%s] DeleteChat TRACE user[%d] deleted chat [%d]\n", reqId, user.Id, id)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("[%s] DeleteChat TRACE distributes user[%d] closes chat[%d]\n", reqId, user.Id, chat.Id)
		err = handler.DistributeChat(app, chat, user, nil, user, event.ChatClose)
		if err != nil {
			log.Printf("[%s] DeleteChat ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		log.Printf("[%s] DeleteChat TRACE distributes user[%d] deletes chat[%d]\n", reqId, user.Id, chat.Id)
		err = handler.DistributeChat(app, chat, user, nil, user, event.ChatDrop)
		if err != nil {
			log.Printf("[%s] DeleteChat ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
		log.Printf("[%s] DeleteChat TRACE distributes user[%d] expel all from chat[%d]\n", reqId, user.Id, chat.Id)
		err = handler.DistributeChat(app, chat, user, nil, nil, event.ChatExpel)
		if err != nil {
			log.Printf("[%s] ExpelUser ERROR cannot distribute chat user expel, %s\n", reqId, err)
			return
		}
	}()
	go func(chatId int, userId uint) {
		defer wg.Done()
		log.Printf("[%s] DeleteChat TRACE user[%d] deletes chat[%d]\n", reqId, userId, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[DELETED_C]"))
	}(chat.Id, user.Id)
	wg.Wait()
}
