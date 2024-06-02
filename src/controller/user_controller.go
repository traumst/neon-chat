package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	d "prplchat/src/db"
	"prplchat/src/handler"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/model/template"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func InviteUser(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
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
	chatId, err := strconv.Atoi(r.FormValue("chatId"))
	if err != nil {
		log.Printf("[%s] InviteUser ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Chat not found"))
		return
	}
	inviteeName := r.FormValue("invitee")
	inviteeName = utils.ReplaceWithSingleSpace(inviteeName)
	inviteeName = utils.RemoveSpecialChars(inviteeName)
	if len(inviteeName) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad invitee name"))
		return
	}
	invitee, err := db.SearchUser(inviteeName)
	if err != nil || invitee == nil {
		log.Printf("[%s] InviteUser ERROR invitee not found, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invitee not found [%s]", inviteeName)))
		return
	}
	log.Printf("[%s] InviteUser TRACE inviting[%d] to chat[%d]\n", reqId, invitee.Id, chatId)
	appInvitee := handler.UserFromDB(*invitee)
	err = app.InviteUser(user.Id, chatId, &appInvitee)
	if err != nil {
		log.Printf("[%s] InviteUser ERROR invite, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to invite user [%s]", invitee.Name)))
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil || chat == nil {
		log.Printf("[%s] InviteUser ERROR user[%d] cannot invite into chat[%d], %s\n", reqId, user.Id, chatId, err.Error())
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
			log.Printf("[%s] InviteUser WARN cannot distribute chat invite, %s\n", reqId, err.Error())
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
			log.Printf("[%s] InviteUser ERROR cannot template user[%d], %s\n", reqId, chatId, err.Error())
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

func ExpelUser(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
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
		// &template.InfoMessage{
		// 	Header: "User is not authenticated",
		// 	Body:   "Your session has probably expired",
		// 	Footer: "Reload the page and try again",
		// }
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelledUserId := r.FormValue("userid")
	expelledId, err := strconv.Atoi(expelledUserId)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR expelled id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbExpelled, err := db.GetUser(uint(expelledId))
	//expelled, err := app.GetUser(uint(expelledId))
	if err != nil || dbExpelled == nil {
		log.Printf("[%s] ExpelUser ERROR expelled, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR cannot find chat[%d], %s\n", reqId, chatId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("[%s] ExpelUser TRACE removing[%d] from chat[%d]\n", reqId, dbExpelled.Id, chatId)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(expelled a.User) {
		defer wg.Done()
		log.Printf("[%s] ExpelUser TRACE distributing user[%d] removed[%d] from chat[%d]\n",
			reqId, user.Id, dbExpelled.Id, chat.Id)
		err := handler.DistributeChat(app, chat, user, &expelled, &expelled, event.ChatClose)
		if err != nil {
			log.Printf("[%s] ExpelUser ERROR cannot distribute chat close, %s\n", reqId, err.Error())
			return
		}
		err = handler.DistributeChat(app, chat, user, &expelled, &expelled, event.ChatDrop)
		if err != nil {
			log.Printf("[%s] ExpelUser ERROR cannot distribute chat deleted, %s\n", reqId, err.Error())
			return
		}
		err = handler.DistributeChat(app, chat, user, nil, &expelled, event.ChatExpel)
		if err != nil {
			log.Printf("[%s] ExpelUser ERROR cannot distribute chat user expel, %s\n", reqId, err.Error())
			return
		}
	}(handler.UserFromDB(*dbExpelled))

	go func() {
		defer wg.Done()
		log.Printf("[%s] ExpelUser TRACE user[%d] removed[%d] from chat[%d]\n", reqId, user.Id, dbExpelled.Id, chat.Id)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf("~<s>%s</s>~", dbExpelled.Name)))
	}()

	wg.Wait()

	err = app.DropUser(user.Id, chatId, dbExpelled.Id)
	if err != nil {
		log.Printf("[%s] ExpelUser ERROR invite, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("[%s] ExpelUser TRACE chat[%d] owner[%d] removed[%d]\n", reqId, chatId, user.Id, dbExpelled.Id)
}

func LeaveChat(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
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
		// &template.InfoMessage{
		// 	Header: "User is not authenticated",
		// 	Body:   "Your session has probably expired",
		// 	Footer: "Reload the page and try again",
		// }
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatid"))
	if err != nil {
		log.Printf("[%s] LeaveChat ERROR chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
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

	err = app.DropUser(user.Id, chat.Id, user.Id)
	if err != nil {
		log.Printf("[%s] LeaveChat out ERROR dropUser, %s\n", reqId, err.Error())
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
			log.Printf("[%s] LeaveChat ERROR cannot distribute chat close, %s\n", reqId, err.Error())
			return
		}
		log.Printf("[%s] LeaveChat TRACE distributed chat close", reqId)
		err = handler.DistributeChat(app, chat, user, user, user, event.ChatDrop)
		if err != nil {
			log.Printf("[%s] LeaveChat ERROR cannot distribute chat deleted, %s\n", reqId, err.Error())
			return
		}
		log.Printf("[%s] LeaveChat TRACE distributed chat deleted", reqId)
		err = handler.DistributeChat(app, chat, user, nil, user, event.ChatLeave)
		if err != nil {
			log.Printf("[%s] LeaveChat ERROR cannot distribute chat user drop, %s\n", reqId, err.Error())
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

func ChangeUser(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] ChangeUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] ChangeUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("bad verb"))
		return
	}
	log.Printf("[%s] ChangeUser TRACE check login\n", reqId)
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] ChangeUser WARN unauthenticated, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unauthenticated"))
		return
	}
	newName := r.FormValue("new-user-name")
	if newName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user did not change"))
		return
	}
	err = db.UpdateUserName(user.Id, newName)
	if err != nil {
		log.Printf("[%s] ChangeUser ERROR failed to update user[%d] in db, %s\n", h.GetReqId(r), user.Id, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user update failed"))
		return
	}
	err = handler.DistributeUserChange(app, user, event.UserChange)
	if err != nil {
		log.Printf("[%s] ChangeUser ERROR failed to distribute user change, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[partial]"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[ok]"))
}
