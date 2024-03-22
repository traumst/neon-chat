package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go.chat/handler"
	"go.chat/model"
	e "go.chat/model/event"
	"go.chat/model/template"
	"go.chat/utils"
)

// TODO if try to open missing chat - fire event to remove
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("--%s-> OpenChat TRACE html template\n", reqId)
	html, err := openChat.Template(user).HTML()
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
	if chatName == "" {
		log.Printf("<-%s-- AddChat ERROR chat name [%s]\n", reqId, chatName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
	log.Printf("--%s-> AddChat TRACE templating chat[%s][%d]\n", reqId, chatName, chatID)

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

func InviteUser(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> InviteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- InviteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := handler.DistributeChat(app, chat, user, invitee, invitee, e.ChatInvite)
		if err != nil {
			log.Printf("<-%s-- InviteUser ERROR cannot distribute chat invite, %s\n", reqId, err)
		}
	}()
	go func() {
		defer wg.Done()
		template := template.MemberTemplate{
			ChatID: chatID,
			Name:   chat.Name,
			User:   invitee,
			Viewer: chat.Owner,
			Owner:  chat.Owner,
		}
		html, err := template.ShortHTML()
		if err != nil {
			log.Printf("<-%s-- InviteUser ERROR cannot template chat[%d], %s\n", reqId, chatID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(html))
	}()
	wg.Wait()

	log.Printf("<-%s-- InviteUser TRACE user [%s] added to chat [%d] by user [%s]\n",
		reqId, invitee, chatID, user)
}

func ExpelUser(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- DeleteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> DeleteUser TRACE check login\n", reqId)
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR auth, %s\n", reqId, err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	chatID, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelled := r.FormValue("userid")
	chat, err := app.GetChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR cannot find chat[%d], %s\n", reqId, chatID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> DeleteUser TRACE removing[%s] from chat[%d]\n", reqId, expelled, chatID)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("--%s-∞ DeleteUser TRACE distributing user[%s] removed[%s] from chat[%d]\n",
			reqId, user, expelled, chat.ID)
		err := handler.DistributeChat(app, chat, user, expelled, expelled, e.ChatClose)
		if err != nil {
			log.Printf("<-%s-- DeleteUser ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, expelled, expelled, e.ChatDeleted)
		if err != nil {
			log.Printf("<-%s-- DeleteUser ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, "", expelled, e.ChatExpel)
		if err != nil {
			log.Printf("<-%s-- DeleteUser ERROR cannot distribute chat user expel, %s\n", reqId, err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- DeleteUser TRACE user[%s] removed[%s] from chat[%d]\n", reqId, user, expelled, chat.ID)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf("expelled <s>%s</s>", expelled)))
	}()
	wg.Wait()

	err = app.DropUser(user, chatID, expelled)
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> DeleteUser TRACE chat[%d] owner[%s] removed[%s]\n", reqId, chatID, user, expelled)
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
	welcome := template.WelcomeTemplate{ActiveUser: user}
	html, err := welcome.HTML()
	if err != nil {
		log.Printf("<-%s-- CloseChat ERROR cannot template welcome page, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("<-%s-- CloseChat TRACE user[%s] closed chat [%d]\n", reqId, user, id)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func LeaveChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> LeaveChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- LeaveChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> LeaveChat TRACE check login\n", reqId)
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("<-%s-- LeaveChat ERROR auth, %s\n", reqId, err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	chatID, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("<-%s-- LeaveChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user, chatID)
	if err != nil {
		log.Printf("<-%s-- LeaveChat ERROR cannot find chat[%d], %s\n", reqId, chatID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> LeaveChat TRACE removing[%s] to chat[%d]\n", reqId, user, chat.ID)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("--%s-∞ LeaveChat TRACE distributing user[%s] left chat[%d]\n", reqId, user, chat.ID)
		err := handler.DistributeChat(app, chat, user, user, user, e.ChatClose)
		if err != nil {
			log.Printf("<-%s-- LeaveChat ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, user, user, e.ChatDeleted)
		if err != nil {
			log.Printf("<-%s-- LeaveChat ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, "", user, e.ChatExpel)
		if err != nil {
			log.Printf("<-%s-- LeaveChat ERROR cannot distribute chat user drop, %s\n", reqId, err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- LeaveChat TRACE user[%s] left chat[%d]\n", reqId, user, chat.ID)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[LEFT_U]"))
	}()
	wg.Wait()

	err = app.DropUser(user, chat.ID, user)
	if err != nil {
		log.Printf("<-%s-- LeaveChat ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> LeaveChat TRACE chat[%d] removed[%s]\n", reqId, chatID, user)
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

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- DeleteChat TRACE user[%s] deletes chat [%d]\n", reqId, user, id)
		err = handler.DistributeChat(app, chat, user, "", "", e.ChatClose)
		if err != nil {
			log.Printf("<-%s-- DeleteChat ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, "", "", e.ChatDeleted)
		if err != nil {
			log.Printf("<-%s-- DeleteChat ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- DeleteUser TRACE user[%s] deletes chat[%d]\n", reqId, user, chat.ID)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[DELETED_C]"))
	}()
	wg.Wait()

	// TODO this needs to move and add recovery
	err = app.DeleteChat(user, chat)
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR remove chat[%d] from [%s], %s\n", reqId, id, chat.Name, err)
	} else {
		log.Printf("<-%s-- DeleteChat TRACE user[%s] deleted chat [%d]\n", reqId, user, id)
	}
}
