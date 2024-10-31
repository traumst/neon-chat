package controller

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src/app"
	"neon-chat/src/consts"
	"neon-chat/src/db"
	"neon-chat/src/event"
	"neon-chat/src/handler/pub"
	"neon-chat/src/sse"
	"neon-chat/src/state"
	h "neon-chat/src/utils/http"
)

func AddAvatar(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] AddAvatar\n", reqId)
	if r.Method != "POST" {
		log.Printf("TRACE [%s] auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	err := r.ParseMultipartForm(consts.MaxUploadBytesSize)
	if err != nil {
		log.Printf("ERROR [%s] multipart failed, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid multipart input"))
		return
	}
	file, info, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("ERROR [%s] reading input file failed, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	defer file.Close()
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	avatar, err := pub.UpdateAvatar(dbConn, user.Id, &file, info)
	if err != nil {
		log.Printf("ERROR  failed to update to avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, "[fail]", http.StatusBadRequest)
		return
	}
	tmpl := avatar.Template(user)
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("ERROR  failed to template avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, fmt.Sprintf("failed to template avatar[%d]", avatar.Id), http.StatusBadRequest)
		return
	}
	state := r.Context().Value(consts.AppState).(*state.State)
	if err = sse.DistributeAvatarChange(state, user, avatar, event.AvatarChange); err != nil {
		log.Printf("ERROR  failed to distribute avatar[%s] update, %s", info.Filename, err.Error())
	}
	w.(*h.StatefulWriter).IndicateChanges()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func GetAvatar(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] GetAvatar\n", reqId)
	if r.Method != "GET" {
		log.Printf("TRACE [%s] auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	avatar, err := pub.GetAvatar(dbConn.Conn, user.Id)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
		return
	}
	tmpl := avatar.Template(user)
	html, err := tmpl.HTML()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to template avatar[%d]", avatar.Id), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
