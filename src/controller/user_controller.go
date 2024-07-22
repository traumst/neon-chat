package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	d "prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func InviteUser(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] InviteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] InviteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] InviteUser WARN user, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	currChatId, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("[%s] InviteUser ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Chat not found"))
		return
	}
	chatId := uint(currChatId)

	inviteeName := r.FormValue("invitee")
	inviteeName = utils.ReplaceWithSingleSpace(inviteeName)
	inviteeName = utils.RemoveSpecialChars(inviteeName)
	if len(inviteeName) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad invitee name"))
		return
	}

	appInvitee, err := handler.FindUser(app, db, inviteeName)
	if err != nil {
		log.Printf("[%s] InviteUser ERROR invitee not found [%s], %s\n", reqId, inviteeName, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Invitee not found [%s]", inviteeName)))
		return
	} else if appInvitee == nil {
		log.Printf("[%s] InviteUser WARN invitee not found [%s]\n", reqId, inviteeName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invitee not found [%s]", inviteeName)))
		return
	}

	chat, err := handler.GetChat(app, db, user, chatId)
	if err != nil {
		log.Printf("[%s] InviteUser ERROR user[%d] cannot invite into chat[%d], %s\n", reqId, user.Id, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Cannot retrieve chat [%d] at this moment", chatId)))
		return
	} else if chat == nil {
		log.Printf("[%s] InviteUser WARN user[%d] cannot invite into chat[%d]\n", reqId, user.Id, chatId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Chat not found [%d] by user [%d]", chatId, user.Id)))
		return
	}

	err = sse.DistributeChat(app, chat, user, appInvitee, appInvitee, event.ChatInvite)
	if err != nil {
		log.Printf("[%s] InviteUser WARN cannot distribute chat invite, %s\n", reqId, err.Error())
	}

	template := template.UserTemplate{
		ChatId:      chatId,
		ChatOwnerId: chat.Owner.Id,
		UserId:      appInvitee.Id,
		UserName:    appInvitee.Name,
		UserEmail:   appInvitee.Email,
		ViewerId:    user.Id,
	}
	html, err := template.HTML()
	if err != nil {
		log.Printf("[%s] InviteUser ERROR cannot template user[%d], %s\n", reqId, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] InviteUser TRACE user[%d] added to chat[%d] by user[%d]\n",
		reqId, appInvitee.Id, chatId, user.Id)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func ExpelUser(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] ExpelUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] ExpelUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] ExpelUser WARN user, %s\n", h.GetReqId(r), err.Error())
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}

	currChatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId := uint(currChatId)

	expelledUserId := r.FormValue("userid")
	expelledId, err := strconv.Atoi(expelledUserId)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR expelled id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appExpelled, err := handler.ExpelUser(app, db, user, chatId, uint(expelledId))
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR failed to expell, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chat, err := handler.GetChat(app, db, user, chatId)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR cannot find chat[%d], %s\n", reqId, chatId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = sse.DistributeChat(app, chat, user, appExpelled, appExpelled, event.ChatClose)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR cannot distribute chat close, %s\n", reqId, err.Error())
		return
	}
	err = sse.DistributeChat(app, chat, user, appExpelled, appExpelled, event.ChatDrop)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR cannot distribute chat deleted, %s\n", reqId, err.Error())
		return
	}
	err = sse.DistributeChat(app, chat, user, nil, appExpelled, event.ChatExpel)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR cannot distribute chat user expel, %s\n", reqId, err.Error())
		return
	}

	log.Printf("[%s] ExpelUser TRACE chat[%d] owner[%d] removed[%d]\n", reqId, chatId, user.Id, appExpelled.Id)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf("~<s>%s</s>~", appExpelled.Name)))
}

func LeaveChat(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] LeaveChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] LeaveChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("[%s] LeaveChat TRACE check login\n", reqId)
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] LeaveChat WARN user, %s\n", h.GetReqId(r), err.Error())
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	currChatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatId := uint(currChatId)

	chat, err := handler.GetChat(app, db, user, chatId)
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR cannot find chat[%d], %s\n", reqId, chatId, err.Error())
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

	// TODO dedicated method?
	expelled, err := handler.ExpelUser(app, db, user, chatId, user.Id)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR failed to expell, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("[%s] LeaveChat TRACE distributing user[%d] left chat[%d]\n", reqId, expelled.Id, chat.Id)
	err = sse.DistributeChat(app, chat, expelled, expelled, expelled, event.ChatClose)
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR cannot distribute chat close, %s\n", reqId, err.Error())
		return
	}
	log.Printf("[%s] LeaveChat TRACE distributed chat close", reqId)
	err = sse.DistributeChat(app, chat, expelled, expelled, expelled, event.ChatDrop)
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR cannot distribute chat deleted, %s\n", reqId, err.Error())
		return
	}
	log.Printf("[%s] LeaveChat TRACE distributed chat deleted", reqId)
	err = sse.DistributeChat(app, chat, expelled, nil, expelled, event.ChatLeave)
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR cannot distribute chat user drop, %s\n", reqId, err.Error())
		return
	}
	log.Printf("[%s] LeaveChat TRACE distributed chat leave", reqId)

	log.Printf("[%s] LeaveChat TRACE user[%d] left chat[%d]\n", reqId, expelled.Id, chat.Id)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("[LEFT_U]"))
}

func ChangeUser(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] ChangeUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] ChangeUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("[verb]"))
		return
	}
	log.Printf("[%s] ChangeUser TRACE check login\n", reqId)
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] ChangeUser WARN unauthenticated, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[auth]"))
		return
	}
	newName := r.FormValue("new-user-name")
	log.Printf("[%s] ChangeUser TRACE new name: %s\n", reqId, newName)
	if newName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[noop]"))
		return
	}
	user.Name = newName

	updatedUser, err := handler.UpdateUser(app, db, user)
	if err != nil {
		log.Printf("[%s] ChangeUser ERROR failed to update user[%d], %s\n", h.GetReqId(r), user.Id, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[fail]"))
		return
	}

	err = sse.DistributeUserChange(app, nil, updatedUser, event.UserChange)
	if err != nil {
		log.Printf("[%s] ChangeUser ERROR failed to distribute user change, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[partial]"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[ok]"))
}

func SearchUsers(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] SearchUsers TRACE IN\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] SearchUsers TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("[verb]"))
		return
	}
	log.Printf("[%s] SearchUsers TRACE check login\n", reqId)
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] SearchUsers WARN unauthenticated, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[auth]"))
		return
	}
	name := r.FormValue("invitee")
	log.Printf("[%s] ChangeUser TRACE new name: %s\n", reqId, name)
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[noop]"))
		return
	}
	users, err := handler.FindUsers(db, name)
	if err != nil {
		log.Printf("[%s] SearchUsers INFO no users matching[%s], %s\n", h.GetReqId(r), name, err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[NoMatch]"))
	}

	html := ""
	for _, appUser := range users {
		tmpl := appUser.Template(0, 0, appUser.Id)
		option, err := tmpl.ShortHTML()
		if err != nil {
			log.Printf("[%s] SearchUsers ERROR failed to template user[%d], %s\n",
				h.GetReqId(r), appUser.Id, err.Error())
			continue
		}
		if len(option) == 0 {
			log.Printf("[%s] SearchUsers ERROR user[%d] has no option\n", h.GetReqId(r), appUser.Id)
			continue
		}
		html += fmt.Sprintf("%s\n%s", html, option)
	}

	if len(html) == 0 {
		log.Printf("[%s] SearchUsers ERROR empty response for users matching[%s]\n", h.GetReqId(r), name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed mapping options"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
