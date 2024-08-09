package controller

import (
	"log"
	"net/http"

	"prplchat/src/convert"
	d "prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	t "prplchat/src/model/template"
	h "prplchat/src/utils/http"
)

func OpenSettings(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] OpenSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] OpenSettings TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] OpenSettings WARN user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User is not authorized"))
		return
	}

	openChatId := app.GetOpenChat(user.Id)
	chatOwnerId := uint(0)
	if openChatId > 0 {
		chat, err := db.GetChat(openChatId)
		if err != nil || chat == nil {
			log.Printf("[%s] OpenSettings ERROR retrieving open chat[%d] data %s\n", reqId, openChatId, err)
		} else {
			chatOwnerId = chat.OwnerId
		}
	}

	var avatarTmpl *t.AvatarTemplate
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

func CloseSettings(app *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] CloseSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] CloseSettings TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("User is unauthorized"))
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] CloseSettings WARN user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("User is unauthorized"))
		return
	}
	var html string
	openChatId := app.GetOpenChat(user.Id)
	openChat, err := db.GetChat(openChatId)
	if openChat != nil {
		chat := convert.ChatDBToApp(openChat)
		dbUsers, err := db.GetChatUsers(user.Id)
		if err != nil {
			log.Printf("[%s] CloseSettings ERROR getting chat[%d] users %s\n", h.GetReqId(r), openChatId, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get chat users"))
			return
		}
		var chatUsers []*a.User
		for _, dbUser := range dbUsers {
			chatUsers = append(chatUsers, convert.UserDBToApp(&dbUser))
		}
		html, err = chat.Template(user, user, chatUsers).HTML()
	} else {
		welcome := t.WelcomeTemplate{User: *user.Template(0, 0, 0)}
		html, err = welcome.HTML()
	}
	if err != nil {
		log.Printf("[%s] CloseSettings ERROR  %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
