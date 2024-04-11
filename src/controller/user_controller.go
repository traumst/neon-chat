package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go.chat/src/db"
	"go.chat/src/handler"
	"go.chat/src/model"
	e "go.chat/src/model/event"
	"go.chat/src/model/template"
	"go.chat/src/utils"
)

func InviteUser(app *model.AppState, conn *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> InviteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- InviteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("--%s-> InviteUser WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	chatId, err := strconv.Atoi(r.FormValue("chatId"))
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
	log.Printf("--%s-> InviteUser TRACE inviting[%d] to chat[%d]\n", reqId, invitee.Id, chatId)
	err = app.InviteUser(user.Id, chatId, invitee)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR invite, %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chat, err := app.GetChat(user.Id, chatId)
	if err != nil {
		log.Printf("<-%s-- InviteUser ERROR cannot find chat[%d], %s\n", reqId, chatId, err)
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
			ChatId:   chatId,
			ChatName: chat.Name,
			Viewer:   template.UserTemplate{Id: chat.Owner.Id, Name: chat.Owner.Name},
			Owner:    template.UserTemplate{Id: chat.Owner.Id, Name: chat.Owner.Name},
		}
		html, err := template.ShortHTML()
		if err != nil {
			log.Printf("<-%s-- InviteUser ERROR cannot template chat[%d], %s\n", reqId, chatId, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(html))
	}()
	wg.Wait()

	log.Printf("<-%s-- InviteUser TRACE user[%d] added to chat[%d] by user[%d]\n",
		reqId, invitee.Id, chatId, user.Id)
}

func ExpelUser(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> DeleteUser\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- DeleteUser TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
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

func LeaveChat(app *model.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := utils.GetReqId(r)
	log.Printf("--%s-> LeaveChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- LeaveChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("--%s-> LeaveChat TRACE check login\n", reqId)
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("--%s-> LeaveChat WARN user, %s\n", utils.GetReqId(r), err)
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
