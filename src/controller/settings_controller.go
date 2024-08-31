package controller

import (
	"log"
	"net/http"

	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/state"
	t "neon-chat/src/model/template"
	h "neon-chat/src/utils/http"
)

func OpenSettings(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] OpenSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] OpenSettings TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] OpenSettings WARN user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User is not authorized"))
		return
	}

	openChatId := state.GetOpenChat(user.Id)
	chatOwnerId := uint(0)
	if openChatId > 0 {
		chat, err := db.GetChat(openChatId)
		if err != nil || chat == nil {
			log.Printf("[%s] OpenSettings ERROR retrieving open chat[%d] data %s\n", reqId, openChatId, err)
		} else {
			chatOwnerId = chat.OwnerId
		}
	}

	var avatarTmpl t.AvatarTemplate
	if avatar, _ := db.GetAvatar(user.Id); avatar != nil {
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

func CloseSettings(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] CloseSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] CloseSettings TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("User is unauthorized"))
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] CloseSettings WARN user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("User is unauthorized"))
		return
	}
	var html string
	openChat := handler.TemplateOpenChat(state, db, user)
	if openChat == nil {
		html, err = handler.TemplateWelcome(user)
	} else {
		html, err = openChat.HTML()
	}
	if err != nil {
		log.Printf("[%s] CloseSettings ERROR failed to template, chat[%t] %s\n", h.GetReqId(r), openChat == nil, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
