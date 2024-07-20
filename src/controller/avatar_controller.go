package controller

import (
	"fmt"
	"log"
	"net/http"

	"prplchat/src/db"
	"prplchat/src/handler"
	"prplchat/src/handler/sse"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func AddAvatar(app *state.State, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] AddAvatar\n", reqId)
	if r.Method != "POST" {
		log.Printf("[%s] AddAvatar TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
	if user == nil {
		log.Printf("[%s] AddAvatar INFO user is not authorized, %s\n", h.GetReqId(r), err)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	err = r.ParseMultipartForm(handler.MaxUploadSize)
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
	avatar := handler.AvatarFromDB(*saved)
	tmpl := avatar.Template(user)
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("controller.AddAvatar ERROR failed to template avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, fmt.Sprintf("failed to template avatar[%d]", avatar.Id), http.StatusBadRequest)
		return
	}
	if err = sse.DistributeAvatarChange(app, user, &avatar, event.AvatarChange); err != nil {
		log.Printf("controller.AddAvatar ERROR failed to distribute avatar[%s] update, %s", info.Filename, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func GetAvatar(app *state.State, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("[%s] GetAvatar\n", reqId)
	if r.Method != "GET" {
		log.Printf("[%s] GetAvatar TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(app, db, w, r)
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
