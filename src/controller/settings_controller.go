package controller

import (
	"log"
	"net/http"

	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/controller/shared"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/state"
	"neon-chat/src/template"
)

func OpenSettings(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] '%s' '%s'\n", reqId, r.Method, r.RequestURI)
	if r.Method != "GET" {
		log.Printf("TRACE [%s] '%s' is not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	openChatId := state.GetOpenChat(user.Id)
	chatOwnerId := uint(0)
	if openChatId > 0 {
		chat, err := db.GetChat(dbConn.Conn, openChatId)
		if err != nil || chat == nil {
			log.Printf("ERROR [%s] retrieving open chat[%d] data %s\n", reqId, openChatId, err)
		} else {
			chatOwnerId = chat.OwnerId
		}
	}
	var avatarTmpl template.AvatarTemplate
	if avatar, _ := db.GetAvatar(dbConn.Conn, user.Id); avatar != nil {
		appAvatar := convert.AvatarDBToApp(avatar)
		avatarTmpl = appAvatar.Template(user)
	}
	settings := template.UserSettingsTemplate{
		ChatId:      openChatId,
		ChatOwnerId: chatOwnerId,
		UserId:      user.Id,
		UserName:    user.Name,
		ViewerId:    user.Id,
		Avatar:      avatarTmpl,
	}
	html, err := settings.HTML()
	if err != nil {
		log.Printf("ERROR [%s] failed to template settings: %s\n", reqId, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func CloseSettings(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] '%s' '%s'\n", reqId, r.Method, r.RequestURI)
	if r.Method != "GET" {
		log.Printf("TRACE [%s] '%s' is not allowed at '%s'\n", reqId, r.Method, r.RequestURI)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User is unauthorized"))
		return
	}
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	state := r.Context().Value(consts.AppState).(*state.State)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	var html string
	var err error
	openChat := shared.TemplateOpenChat(state, dbConn, user)
	if openChat == nil {
		html, err = shared.TemplateWelcome(user)
	} else {
		html, err = openChat.HTML()
	}
	if err != nil {
		log.Printf("ERROR [%s] failed to template open chat, wasNull[%t], %s\n", reqId, openChat == nil, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
