package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go.chat/src/db"
	"go.chat/src/handler"
	"go.chat/src/model/event"
	"go.chat/src/model/template"
	"go.chat/src/utils"
	h "go.chat/src/utils/http"
)

func InviteUser(app *handler.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] InviteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] InviteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] InviteUser WARN user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("[%s] InviteUser ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Chat not found"))
		return
	}
	inviteeName := r.FormValue("invitee")
	inviteeName = utils.TrimSpaces(inviteeName)
	inviteeName = utils.TrimSpecial(inviteeName)
	if inviteeName == "" || len(inviteeName) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad invitee name"))
		return
	}
	invitee, err := conn.SearchUser(inviteeName)
	if err != nil || invitee == nil {
		log.Printf("[%s] InviteUser ERROR invitee not found, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invitee not found [%s]", inviteeName)))
		return
	}
	log.Printf("[%s] InviteUser TRACE inviting[%d] to chat[%d]\n", reqId, invitee.Id, chatId)
	appInvitee := handler.UserFromDB(*invitee)
	err = app.InviteUser(user.Id, chatId, &appInvitee)
	if err != nil {
		log.Printf("[%s] InviteUser ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to invite user [%s]", invitee.Name)))
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil || chat == nil {
		log.Printf("[%s] InviteUser ERROR user[%d] cannot invite into chat[%d], %s\n", reqId, user.Id, chatId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Cannot invite user [%s] into this chat", invitee.Name)))
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := handler.DistributeChat(app, chat, user, &appInvitee, &appInvitee, event.ChatInvite)
		if err != nil {
			log.Printf("[%s] InviteUser WARN cannot distribute chat invite, %s\n", reqId, err)
		}
	}()
	go func() {
		defer wg.Done()
		template := template.UserTemplate{
			ChatId:      chatId,
			ChatOwnerId: chat.Owner.Id,
			UserId:      invitee.Id,
			UserName:    invitee.Name,
			UserEmail:   invitee.Email,
			ViewerId:    user.Id,
		}
		html, err := template.HTML()
		if err != nil {
			log.Printf("[%s] InviteUser ERROR cannot template user[%d], %s\n", reqId, chatId, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(html))
	}()
	wg.Wait()

	log.Printf("[%s] InviteUser TRACE user[%d] added to chat[%d] by user[%d]\n",
		reqId, invitee.Id, chatId, user.Id)
}

func ExpelUser(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] ExpelUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] ExpelUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] ExpelUser WARN user, %s\n", h.GetReqId(r), err)
		RenderLogin(w, r, &template.InfoMessage{
			Header: "User is not authenticated",
			Body:   "Your session has probably expired",
			Footer: "Reload the page and try again",
		})
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelledUserId := r.FormValue("userid")
	expelledId, err := strconv.Atoi(expelledUserId)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR expelled id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelled, err := app.GetUser(uint(expelledId))
	if err != nil || expelled == nil {
		log.Printf("[%s] ExpelUser ERROR expelled, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR cannot find chat[%d], %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("[%s] ExpelUser TRACE removing[%d] from chat[%d]\n", reqId, expelled.Id, chatId)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("[%s] ExpelUser TRACE distributing user[%d] removed[%d] from chat[%d]\n",
			reqId, user.Id, expelled.Id, chat.Id)
		err := handler.DistributeChat(app, chat, user, expelled, expelled, event.ChatClose)
		if err != nil {
			log.Printf("[%s] ExpelUser ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, expelled, expelled, event.ChatDrop)
		if err != nil {
			log.Printf("[%s] ExpelUser ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
		err = handler.DistributeChat(app, chat, user, nil, expelled, event.ChatExpel)
		if err != nil {
			log.Printf("[%s] ExpelUser ERROR cannot distribute chat user expel, %s\n", reqId, err)
			return
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("[%s] ExpelUser TRACE user[%d] removed[%d] from chat[%d]\n", reqId, user.Id, expelled.Id, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf("~<s>%s</s>~", expelled.Name)))
	}()
	wg.Wait()

	err = app.DropUser(user.Id, chatId, expelled.Id)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("[%s] ExpelUser TRACE chat[%d] owner[%d] removed[%d]\n", reqId, chatId, user.Id, expelled.Id)
}

func LeaveChat(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] LeaveChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] LeaveChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("[%s] LeaveChat TRACE check login\n", reqId)
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] LeaveChat WARN user, %s\n", h.GetReqId(r), err)
		RenderLogin(w, r, &template.InfoMessage{
			Header: "User is not authenticated",
			Body:   "Your session has probably expired",
			Footer: "Reload the page and try again",
		})
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR chat id, %s\n", reqId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR cannot find chat[%d], %s\n", reqId, chatId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("[%s] LeaveChat TRACE removing[%d] from chat[%d]\n", reqId, user.Id, chat.Id)
	if user.Id == chat.Owner.Id {
		log.Printf("[%s] LeaveChat ERROR cannot leave chat[%d] as owner\n", reqId, chatId)
		w.Write([]byte("creator cannot leave chat"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.DropUser(user.Id, chat.Id, user.Id)
	if err != nil {
		log.Printf("[%s] LeaveChat out ERROR dropUser, %s\n", reqId, err)
	} else {
		log.Printf("[%s] LeaveChat out TRACE chat[%d] removed[%d]\n", reqId, chatId, user.Id)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("[%s] LeaveChat TRACE distributing user[%d] left chat[%d]\n", reqId, user.Id, chat.Id)
		err := handler.DistributeChat(app, chat, user, user, user, event.ChatClose)
		if err != nil {
			log.Printf("[%s] LeaveChat ERROR cannot distribute chat close, %s\n", reqId, err)
			return
		}
		log.Printf("[%s] LeaveChat TRACE distributed chat close", reqId)
		err = handler.DistributeChat(app, chat, user, user, user, event.ChatDrop)
		if err != nil {
			log.Printf("[%s] LeaveChat ERROR cannot distribute chat deleted, %s\n", reqId, err)
			return
		}
		log.Printf("[%s] LeaveChat TRACE distributed chat deleted", reqId)
		err = handler.DistributeChat(app, chat, user, nil, user, event.ChatLeave)
		if err != nil {
			log.Printf("[%s] LeaveChat ERROR cannot distribute chat user drop, %s\n", reqId, err)
			return
		}
		log.Printf("[%s] LeaveChat TRACE distributed chat leave", reqId)
	}()
	go func() {
		defer wg.Done()
		log.Printf("[%s] LeaveChat TRACE user[%d] left chat[%d]\n", reqId, user.Id, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("[LEFT_U]"))
	}()
	wg.Wait()
}

func ChangeUser(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] ChangeUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] ChangeUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("bad verb"))
		return
	}
	log.Printf("[%s] ChangeUser TRACE check login\n", reqId)
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] ChangeUser WARN unauthenticated, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unauthenticated"))
		return
	}
	newName := r.FormValue("new-user-name")
	if newName == "" {
		log.Printf("[%s] ChangeUser TRACE user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user did not change"))
		return
	}
	err = db.UpdateUserName(user.Id, newName)
	if err != nil {
		log.Printf("[%s] ChangeUser WARN failed to update user[%d], %s\n", h.GetReqId(r), user.Id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user update failed"))
		return
	}
	err = handler.DistributeUserChange(app, user, event.UserChange)
	if err != nil {
		log.Printf("[%s] ChangeUser ERROR failed to distribute user change, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[partial]"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[ok]"))
}
