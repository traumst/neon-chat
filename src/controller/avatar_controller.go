package controller

import (
	"fmt"
	"log"
	"net/http"

	"neon-chat/src/convert"
	d "neon-chat/src/db"
	"neon-chat/src/handler"
	"neon-chat/src/handler/sse"
	"neon-chat/src/handler/state"
	a "neon-chat/src/model/app"
	"neon-chat/src/model/event"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
)

func AddAvatar(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] AddAvatar\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] AddAvatar TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if user == nil {
		log.Printf("[%s] AddAvatar INFO user is not authorized, %s\n", h.GetReqId(r), err)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	err = r.ParseMultipartForm(utils.MaxUploadBytesSize)
	if err != nil {
		log.Printf("[%s] AddAvatar ERROR multipart failed, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid multipart input"))
		return
	}
	file, info, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("[%s] AddAvatar ERROR reading input file failed, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	defer file.Close()
	saved, err := handler.UpdateAvatar(db, user.Id, &file, info)
	if err != nil {
		log.Printf("controller.AddAvatar ERROR failed to update to avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, "[fail]", http.StatusBadRequest)
		return
	}
	avatar := convert.AvatarDBToApp(saved)
	tmpl := avatar.Template(user)
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("controller.AddAvatar ERROR failed to template avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, fmt.Sprintf("failed to template avatar[%d]", avatar.Id), http.StatusBadRequest)
		return
	}
	if err = sse.DistributeAvatarChange(state, user, avatar, event.AvatarChange); err != nil {
		log.Printf("controller.AddAvatar ERROR failed to distribute avatar[%s] update, %s", info.Filename, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func GetAvatar(state *state.State, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] GetAvatar\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] GetAvatar TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(state, db, w, r)
	if user == nil {
		log.Printf("[%s] GetAvatar INFO user is not authorized, %s\n", h.GetReqId(r), err)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	dbAvatar, err := db.GetAvatar(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
		return
	}
	avatar := &a.Avatar{
		Id:     dbAvatar.Id,
		UserId: dbAvatar.UserId,
		Title:  dbAvatar.Title,
		Size:   fmt.Sprintf("%dKB", dbAvatar.Size/utils.KB),
		Image:  dbAvatar.Image,
		Mime:   dbAvatar.Mime,
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
