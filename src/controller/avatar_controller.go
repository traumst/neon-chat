package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"go.chat/src/handler"
	a "go.chat/src/model/app"
	"go.chat/src/utils"
	h "go.chat/src/utils/http"
)

const MaxUploadSize int64 = 10 * utils.KB

func AddAvatar(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("--%s-> OpenChat\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- OpenChat TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if user == nil {
		log.Printf("--%s-> OpenChat INFO user is not authorized, %s\n", h.GetReqId(r), err)
		RenderHome(app, w, r)
		return
	}
	err = r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("file is too big, limit is %dKB", MaxUploadSize), http.StatusBadRequest)
		return
	}
	file, info, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	if info.Size > MaxUploadSize {
		http.Error(w, "file is too large, limit is "+strconv.Itoa(int(MaxUploadSize))+"KB", http.StatusBadRequest)
		return
	} else if info.Filename == "" {
		http.Error(w, "file lacks name", http.StatusBadRequest)
		return
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to load input file", http.StatusBadRequest)
		return
	}
	fileType := http.DetectContentType(fileBytes)
	if fileType != "image/svg+xml" {
		http.Error(w, "file type is not supported: "+fileType, http.StatusBadRequest)
		return
	}
	avatar := a.UserAvatar{
		UserId: user.Id,
		Title:  info.Filename,
		Mime:   fileType,
		Size:   fmt.Sprintf("%dKB", info.Size/utils.KB),
		Image:  fileBytes,
	}
	saved, err := app.AddAvatar(user.Id, &avatar)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save avatar[%s]", info.Filename), http.StatusBadRequest)
		return
	}
	avatar.Id = saved.Id
	tmpl := avatar.Template(user)
	html, err := tmpl.HTML()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to template avatar[%d]", avatar.Id), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func GetAvatar(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("--%s-> GetAvatar\n", reqId)
	if r.Method != "GET" {
		log.Printf("<-%s-- GetAvatar TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if user == nil {
		log.Printf("--%s-> GetAvatar INFO user is not authorized, %s\n", h.GetReqId(r), err)
		RenderHome(app, w, r)
		return
	}
	avatar, err := app.GetAvatar(user.Id)
	if err != nil {
		http.Error(w, "avatar not found", http.StatusNotFound)
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
