package controller

import (
	"log"
	"net/http"

	"go.chat/src/handler"
	t "go.chat/src/model/template"
	h "go.chat/src/utils/http"
)

func OpenSettings(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] OpenSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] OpenSettings TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] OpenSettings WARN user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User is not authorized"))
		return
	}
	var openChatId int
	var chatOwnerId uint
	openChat := app.GetOpenChat(user.Id)
	if openChat != nil {
		openChatId = openChat.Id
		chatOwnerId = openChat.Owner.Id
	} else {
		openChatId = 0
		chatOwnerId = 0
	}
	var avatarTmpl *t.AvatarTemplate
	if avatar, _ := app.GetAvatar(user.Id); avatar != nil {
		avatarTmpl = avatar.Template(user)
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

func CloseSettings(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] CloseSettings\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] CloseSettings TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("User is unauthorized"))
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if err != nil || user == nil {
		log.Printf("[%s] CloseSettings WARN user, %s\n", h.GetReqId(r), err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("User is unauthorized"))
		return
	}
	var html string
	openChat := app.GetOpenChat(user.Id)
	if openChat != nil {
		html, err = openChat.Template(user, user).HTML()
	} else {
		welcome := t.WelcomeTemplate{User: *user.Template(-1, 0, 0)}
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

func ChangeAvatar(app *handler.AppState, w http.ResponseWriter, r *http.Request) {

}
