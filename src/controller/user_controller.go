package controller

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"neon-chat/src/event"
	"neon-chat/src/handler/parse"
	"neon-chat/src/handler/pub"
	"neon-chat/src/sse"
	"neon-chat/src/state"
	"neon-chat/src/template"
	h "neon-chat/src/utils/http"
)

func InviteUser(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] '%s' is not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("ERROR [%s] chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inviteeName, err := parse.ReadFormValueString(r, "invitee")
	if err != nil {
		log.Printf("ERROR [%s] invitee name, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(inviteeName) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad invitee name"))
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, appInvitee, err := pub.InviteUser(state, dbConn, user, chatId, inviteeName)
	if err != nil {
		log.Printf("ERROR [%s] failed to invite user[%d] into chat[%d], %s\n",
			reqId, appInvitee.Id, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to invite user [%d] into chat [%d]", appInvitee.Id, chatId)))
		return
	} else if appChat == nil || appInvitee == nil {
		log.Printf("WARN [%s] user[%d] not found or chat[%d] not found\n",
			reqId, appInvitee.Id, chatId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	err = sse.DistributeChat(state, dbConn.Tx, appChat, user, appInvitee, appInvitee, event.ChatInvite)
	if err != nil {
		log.Printf("WARN [%s] cannot distribute chat invite, %s\n", reqId, err.Error())
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
		log.Printf("ERROR [%s] cannot template user[%d], %s\n", reqId, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.(*h.StatefulWriter).IndicateChanges()
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(html))
}

func ExpelUser(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] '%s' is not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("ERROR [%s] chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expelledId, err := parse.ReadFormValueUint(r, "userid")
	if err != nil {
		log.Printf("ERROR [%s] expelled id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, appExpelled, err := pub.ExpelUser(state, dbConn, user, chatId, expelledId)
	if err != nil {
		log.Printf("ERROR [%s] failed to expell user[%d] from chat[%d], %s\n",
			reqId, expelledId, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = sse.DistributeChat(state, dbConn.Tx, appChat, user, appExpelled, appExpelled, event.ChatClose)
	if err != nil {
		log.Printf("ERROR  cannot distribute chat close, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, appChat, user, appExpelled, appExpelled, event.ChatDrop)
	if err != nil {
		log.Printf("ERROR  cannot distribute chat deleted, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, appChat, user, nil, appExpelled, event.ChatExpel)
	if err != nil {
		log.Printf("ERROR  cannot distribute chat expel, %s\n", err.Error())
	}
	w.(*h.StatefulWriter).IndicateChanges()
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf("~<s>%s</s>~", appExpelled.Name)))
}

func LeaveChat(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] '%s' is not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatId, err := parse.ReadFormValueUint(r, "chatid")
	if err != nil {
		log.Printf("ERROR [%s] chat id, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	appChat, leftUser, err := pub.LeaveChat(state, dbConn, user, chatId)
	if err != nil {
		log.Printf("ERROR [%s] failed to leave chat[%d], %s\n", reqId, chatId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("TRACE  informing users in chat[%d]\n", appChat.Id)
	err = sse.DistributeChat(state, dbConn.Tx, appChat, leftUser, leftUser, leftUser, event.ChatClose)
	if err != nil {
		log.Printf("ERROR  cannot distribute chat close, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, appChat, leftUser, leftUser, leftUser, event.ChatDrop)
	if err != nil {
		log.Printf("ERROR  cannot distribute chat deleted, %s\n", err.Error())
	}
	err = sse.DistributeChat(state, dbConn.Tx, appChat, leftUser, nil, leftUser, event.ChatLeave)
	if err != nil {
		log.Printf("ERROR  cannot distribute chat user drop, %s\n", err.Error())
	}
	w.(*h.StatefulWriter).IndicateChanges()
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("[LEFT_U]"))
}

func ChangeUser(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] '%s' is not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("[verb]"))
		return
	}
	newName := r.FormValue("new-user-name")
	log.Printf("TRACE [%s] new name: %s\n", reqId, newName)
	if newName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[noop]"))
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	user.Name = newName
	updatedUser, err := pub.UpdateUser(state, dbConn.Tx, user)
	if err != nil {
		log.Printf("ERROR [%s] failed to update user[%d], %s\n", reqId, user.Id, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[fail]"))
		return
	}
	err = sse.DistributeUserChange(state, dbConn.Tx, nil, updatedUser, event.UserChange)
	if err != nil {
		log.Printf("ERROR [%s] failed to distribute user change, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[partial]"))
		return
	}
	w.(*h.StatefulWriter).IndicateChanges()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[ok]"))
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] '%s' is not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("[verb]"))
		return
	}
	log.Printf("TRACE [%s] check login\n", reqId)
	name := r.FormValue("invitee")
	log.Printf("TRACE [%s] new name: %s\n", reqId, name)
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[noop]"))
		return
	}
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	users, err := pub.SearchUsers(dbConn.Conn, name)
	if err != nil {
		log.Printf("INFO [%s] no users matching[%s], %s\n", reqId, name, err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[NoMatch]"))
	}
	html := ""
	for _, appUser := range users {
		tmpl := appUser.Template(0, 0, appUser.Id)
		option, err := tmpl.ShortHTML()
		if err != nil {
			log.Printf("ERROR [%s] failed to template user[%d], %s\n",
				reqId, appUser.Id, err.Error())
			continue
		}
		if len(option) == 0 {
			log.Printf("ERROR [%s] user[%d] has no option\n", reqId, appUser.Id)
			continue
		}
		html += fmt.Sprintf("%s\n%s", html, option)
	}
	if len(html) == 0 {
		log.Printf("ERROR [%s] empty response for users matching[%s]\n", reqId, name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed mapping options"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
