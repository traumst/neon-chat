package controller

import (
	"log"
	"net/http"

	"neon-chat/src/consts"
	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler"
	a "neon-chat/src/model/app"
	t "neon-chat/src/model/template"
	"neon-chat/src/state"
)

func OpenSettings(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] OpenSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] OpenSettings TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value(consts.ActiveUser).(*a.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	db := r.Context().Value(consts.DBConn).(*d.DBConn)
	openChatId := state.GetOpenChat(user.Id)
	chatOwnerId := uint(0)
	if openChatId > 0 {
		chat, err := d.GetChat(db.Conn, openChatId)
		if err != nil || chat == nil {
			log.Printf("[%s] OpenSettings ERROR retrieving open chat[%d] data %s\n", reqId, openChatId, err)
		} else {
			chatOwnerId = chat.OwnerId
		}
	}

	var avatarTmpl t.AvatarTemplate
	if avatar, _ := d.GetAvatar(db.Conn, user.Id); avatar != nil {
		appAvatar := convert.AvatarDBToApp(avatar)
		avatarTmpl = appAvatar.Template(user)
	}
	settings := t.UserSettingsTemplate{
		ChatId:      openChatId,
		ChatOwnerId: chatOwnerId,
		UserId:      user.Id,
		UserName:    user.Name,
		ViewerId:    user.Id,
		Avatar:      avatarTmpl,
	}
	html, err := settings.HTML()
	if err != nil {
		log.Printf("[%s] OpenSettings ERROR %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func CloseSettings(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("[%s] CloseSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] CloseSettings TRACE auth does not allow %s\n", reqId, r.Method)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User is unauthorized"))
		return
	}

	user := r.Context().Value(consts.ActiveUser).(*a.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	db := r.Context().Value(consts.DBConn).(*d.DBConn)

	var html string
	var err error
	openChat := handler.TemplateOpenChat(state, db, user)
	if openChat == nil {
		html, err = handler.TemplateWelcome(user)
	} else {
		html, err = openChat.HTML()
	}
	if err != nil {
		log.Printf("[%s] CloseSettings ERROR failed to template, chat[%t] %s\n", reqId, openChat == nil, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
