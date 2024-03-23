package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go.chat/db"
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
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> Home WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	path := utils.ParseUrlPath(r)
	log.Printf("--%s-> OpenChat, %s\n", reqId, path[2])
	chatId, err := strconv.Atoi(path[2])
	if err != nil {
		log.Printf("<-%s-- OpenChat ERROR id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if chatId < 0 {
		log.Printf("<-%s-- OpenChat ERROR chatID, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> OpenChat TRACE chat[%d]\n", reqId, chatId)
	openChat, err := app.OpenChat(user.Id, chatId)
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
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> AddChat WARN user, %s\n", utils.GetReqId(r), err)
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
	chatID := app.AddChat(user, chatName)
	log.Printf("--%s-> AddChat TRACE user[%d] opening chat[%s][%d]\n", reqId, user.Id, chatName, chatID)
	openChat, err := app.OpenChat(user.Id, chatID)
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

func InviteUser(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> InviteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- InviteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		log.Printf("--%s-> InviteUser WARN cookie, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> InviteUser WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	chatID, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inviteeName := r.FormValue("invitee")
	invitee, err := conn.GetUser(inviteeName)
	if err != nil || invitee == nil {
		log.Printf("<-%s-- InviteUser ERROR invitee not found, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> InviteUser TRACE inviting[%d] to chat[%d]\n", reqId, invitee.Id, chatID)
	err = app.InviteUser(user.Id, chatID, invitee)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chat, err := app.GetChat(user.Id, chatID)
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
			User:   invitee.Name,
			Viewer: chat.Owner.Name,
			Owner:  chat.Owner.Name,
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

	log.Printf("<-%s-- InviteUser TRACE user[%d] added to chat[%d] by user[%d]\n",
		reqId, invitee.Id, chatID, user.Id)
}

func ExpelUser(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- DeleteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		log.Printf("--%s-> DeleteUser WARN cookie\n", reqId)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> DeleteUser WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelledUserId := r.FormValue("userid")
	expelledId, err := strconv.Atoi(expelledUserId)
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR expelled id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelled, err := app.GetUser(uint(expelledId))
	if err != nil || expelled == nil {
		log.Printf("<-%s-- DeleteUser ERROR expelled, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR cannot find chat[%d], %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> DeleteUser TRACE removing[%d] from chat[%d]\n", reqId, expelled.Id, chatId)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("--%s-∞ DeleteUser TRACE distributing user[%d] removed[%d] from chat[%d]\n",
			reqId, user.Id, expelled.Id, chat.Id)
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
		err = handler.DistributeChat(app, chat, user, nil, expelled, e.ChatExpel)
		if err != nil {
			log.Printf("<-%s-- DeleteUser ERROR cannot distribute chat user expel, %s\n", reqId, err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- DeleteUser TRACE user[%d] removed[%d] from chat[%d]\n", reqId, user.Id, expelled.Id, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf("expelled <s>%s</s>", expelled.Name)))
	}()
	wg.Wait()

	err = app.DropUser(user.Id, chatId, expelled.Id)
	if err != nil {
		log.Printf("<-%s-- DeleteUser ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> DeleteUser TRACE chat[%d] owner[%d] removed[%d]\n", reqId, chatId, user.Id, expelled.Id)
}

func CloseChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> CloseChat\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		log.Printf("--%s-> DeleteUser WARN cookie\n", reqId)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> DeleteUser WARN user, %s\n", utils.GetReqId(r), err)
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

func LeaveChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> LeaveChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- LeaveChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> LeaveChat TRACE check login\n", reqId)
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		log.Printf("--%s-> DeleteUser WARN cookie\n", reqId)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> DeleteUser WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("<-%s-- LeaveChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil {
		log.Printf("<-%s-- LeaveChat ERROR cannot find chat[%d], %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> LeaveChat TRACE removing[%d] from chat[%d]\n", reqId, user.Id, chat.Id)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("--%s-∞ LeaveChat TRACE distributing user[%d] left chat[%d]\n", reqId, user.Id, chat.Id)
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
		err = handler.DistributeChat(app, chat, user, nil, user, e.ChatExpel)
		if err != nil {
			log.Printf("<-%s-- LeaveChat ERROR cannot distribute chat user drop, %s\n", reqId, err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- LeaveChat TRACE user[%d] left chat[%d]\n", reqId, user.Id, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[LEFT_U]"))
	}()
	wg.Wait()

	err = app.DropUser(user.Id, chat.Id, user.Id)
	if err != nil {
		log.Printf("<-%s-- LeaveChat ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("--%s-> LeaveChat TRACE chat[%d] removed[%d]\n", reqId, chatId, user.Id)
}

func DeleteChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteChat TRACE\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		log.Printf("--%s-> DeleteUser WARN cookie\n", reqId)
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}
	user, err := app.GetUser(cookie.UserId)
	if err != nil || user == nil {
		log.Printf("--%s-> DeleteUser WARN user, %s\n", utils.GetReqId(r), err)
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
		log.Printf("<-%s-- DeleteChat TRACE user[%d] deletes chat[%d]\n", reqId, user.Id, chat.Id)
		err = handler.DistributeChat(app, chat, user, nil, nil, e.ChatClose)
		if err != nil {
			log.Printf("<-%s-- DeleteChat ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, nil, nil, e.ChatDeleted)
		if err != nil {
			log.Printf("<-%s-- DeleteChat ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("<-%s-- DeleteUser TRACE user[%d] deletes chat[%d]\n", reqId, user.Id, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[DELETED_C]"))
	}()
	wg.Wait()

	// TODO this needs to move and add recovery
	err = app.DeleteChat(user.Id, chat)
	if err != nil {
		log.Printf("<-%s-- DeleteChat ERROR remove chat[%d] from [%s], %s\n", reqId, id, chat.Name, err)
	} else {
		log.Printf("<-%s-- DeleteChat TRACE user[%d] deleted chat [%d]\n", reqId, user.Id, id)
	}
}
