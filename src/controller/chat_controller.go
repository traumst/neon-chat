package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"go.chat/src/handler"
	"go.chat/src/model"
	e "go.chat/src/model/event"
	"go.chat/src/model/template"
	"go.chat/src/utils"
)

// TODO if try to open missing chat - fire event to remove
func OpenChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> OpenChat\n", reqId)
	if r.Method != "GET" {
		log.Printf("<-%s-- OpenChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if user == nil {
		log.Printf("--%s-> OpenChat INFO user is not authorized, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	log.Printf("--%s-> OpenChat, %s\n", reqId, path[2])
	chatId, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid chat id %s", path[2])))
		return
	}
	if chatId < 0 {
		log.Printf("<-%s-- OpenChat ERROR chatId, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid chat id %d", chatId)))
		return
	}
	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", reqId, chatId)
	openChat, err := app.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR chat, %s\n", reqId, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Chat not found"))
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", reqId)
	html, err := openChat.Template(user).HTML()
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR html template, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template chat"))
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
	user, err := handler.ReadSession(app, w, r)
	if user == nil {
		log.Printf("--%s-> AddChat INFO user is not authorized, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	chatName := r.FormValue("chatName")
	if chatName == "" {
		log.Printf("<-%s-- AddChat ERROR chat name [%s]\n", reqId, chatName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> AddChat TRACE adding user[%d] chat[%s]\n", reqId, user.Id, chatName)
	chatId := app.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE user[%d] opening chat[%s][%d]\n", reqId, user.Id, chatName, chatId)
	openChat, err := app.OpenChat(user.Id, chatId)
	if err != nil {
		log.Printf("<-%s-- AddChat ERROR chat, %s\n", reqId, err)
		errMsg := fmt.Sprintf("ERROR: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}
	log.Printf("--%s-> AddChat TRACE templating chat[%s][%d]\n", reqId, chatName, chatId)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = handler.DistributeChat(app, openChat, user, user, user, e.ChatCreated)
		if err != nil {
			log.Printf("<-%s-- AddChat ERROR cannot distribute chat header, %s\n", reqId, err)
		}
	}()
	go func() {
		defer wg.Done()
		template := openChat.Template(user)
		html, err := template.HTML()
		if err != nil {
			log.Printf("<--%s-- sendChatContent ERROR cannot template chat [%+v], %s", reqId, template, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("<-%s-- sendChatContent TRACE writing response\n", reqId)
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(html))
	}()
	wg.Wait()
}

func CloseChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> CloseChat\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("--%s-> CloseChat WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	chatIdStr := r.PostFormValue("chatid")
	if chatIdStr == "" {
		log.Printf("<-%s-- CloseChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId, err := strconv.Atoi(chatIdStr)
	if err != nil {
		log.Printf("<-%s-- CloseChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.CloseChat(user.Id, chatId)
	if err != nil {
		log.Printf("<-%s-- CloseChat ERROR close chat[%d] for user[%d], %s\n",
			reqId, chatId, user.Id, err)
	}
	welcome := template.WelcomeTemplate{ActiveUser: user.Name}
	html, err := welcome.HTML()
	if err != nil {
		log.Printf("<-%s-- CloseChat ERROR cannot template welcome page, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- CloseChat TRACE user[%d] closed chat[%d]\n", reqId, user.Id, chatId)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func DeleteChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteChat TRACE\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("--%s-> DeleteChat WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	chatId := r.PostFormValue("chatid")
	if chatId == "" {
		log.Printf("<-%s-- DeleteChat ERROR parse args, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chatId)
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, id)
	if err != nil || chat == nil {
		log.Printf("<-%s-- DeleteChat ERROR cannot get chat[%d] for user[%d]\n", reqId, id, user.Id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- DeleteChat TRACE distributes user[%d] closes chat[%d]\n", reqId, user.Id, chat.Id)
		err = handler.DistributeChat(app, chat, user, nil, nil, e.ChatClose)
		if err != nil {
			log.Printf("<-%s-- DeleteChat ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		log.Printf("<-%s-- DeleteChat TRACE distributes user[%d] deletes chat[%d]\n", reqId, user.Id, chat.Id)
		err = handler.DistributeChat(app, chat, user, nil, nil, e.ChatDeleted)
		if err != nil {
			log.Printf("<-%s-- DeleteChat ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
	}()
	go func(chatId int, userId uint) {
		defer wg.Done()
		log.Printf("<-%s-- DeleteChat TRACE user[%d] deletes chat[%d]\n", reqId, userId, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[DELETED_C]"))
	}(chat.Id, user.Id)
	wg.Wait()

	// TODO this needs to move and add recovery
	err = app.DeleteChat(user.Id, chat)
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR remove chat[%d] from [%s], %s\n", reqId, id, chat.Name, err)
	} else {
		log.Printf("<-%s-- DeleteChat TRACE user[%d] deleted chat [%d]\n", reqId, user.Id, id)
	}
}
