package controller

import (
	"fmt"
	"log"
	"net/http"

	d "prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
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

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] InviteUser ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inviteeName, err := handler.FormValueString(r, "invitee")
	if err != nil {
		log.Printf("[%s] InviteUser ERROR invitee name, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(inviteeName) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad invitee name"))
		return
	}

	appChat, appInvitee, err := handler.HandleUserInvite(app, db, user, chatId, inviteeName)
	if err != nil {
		log.Printf("[%s] InviteUser ERROR failed to invite user[%d] into chat[%d], %s\n",
			reqId, appInvitee.Id, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to invite user [%d] into chat [%d]", appInvitee.Id, chatId)))
		return
	} else if appChat == nil || appInvitee == nil {
		log.Printf("[%s] InviteUser WARN user[%d] not found or chat[%d] not found\n",
			reqId, appInvitee.Id, chatId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	template := template.UserTemplate{
		ChatId:      chatId,
		ChatOwnerId: appChat.OwnerId,
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

	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelledId, err := handler.FormValueUint(r, "userid")
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR expelled id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appExpelled, err := handler.HandleUserExpelled(app, db, user, chatId, expelledId)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR failed to expell user[%d] from chat[%d], %s\n",
			reqId, expelledId, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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
	chatId, err := handler.FormValueUint(r, "chatid")
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.HandleUserLeaveChat(app, db, user, chatId)
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR failed to leave chat[%d], %s\n", reqId, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] LeaveChat TRACE user[%d] left chat[%d]\n", reqId, user.Id, chatId)
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
