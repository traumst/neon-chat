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

var allowedImageFormats = []string{
	"image/svg+xml",
	"image/jpeg",
	"image/gif",
}

func isAllowedImageFormat(mime string) bool {
	for _, allowed := range allowedImageFormats {
		if allowed == mime {
			return true
		}
	}
	return false
}

func AddAvatar(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	reqId := h.GetReqId(r)
	log.Printf("--%s-> AddAvatar\n", reqId)
	if r.Method != "POST" {
		log.Printf("<-%s-- AddAvatar TRACE auth does not allow %s\n", reqId, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET method is allowed"))
		return
	}
	user, err := handler.ReadSession(app, w, r)
	if user == nil {
		log.Printf("--%s-> AddAvatar INFO user is not authorized, %s\n", h.GetReqId(r), err)
		RenderHome(app, w, r)
		return
	}
	err = r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		log.Printf("<-%s-- AddAvatar ERROR multipart failed, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid multipart input"))
		return
	}
	file, info, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("<-%s-- AddAvatar ERROR reading input file failed, %s\n", reqId, err.Error())
		log.Printf("<-%s-- AddAvatar TRACE reading input file failed, %+v\n", reqId, info.Filename)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid input"))
		return
	}
	defer file.Close()
	if info.Size > MaxUploadSize {
		http.Error(w, "file too large "+strconv.Itoa(int(info.Size))+
			", limit is "+strconv.Itoa(int(MaxUploadSize)), http.StatusBadRequest)
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
	if !isAllowedImageFormat(fileType) {
		http.Error(w, "file type is not supported: "+fileType, http.StatusBadRequest)
		return
	}
	avatar := a.Avatar{
		UserId: user.Id,
		Title:  info.Filename,
		Mime:   fileType,
		Size:   fmt.Sprintf("%dKB", info.Size/utils.KB),
		Image:  fileBytes,
	}
	saved, err := app.AddAvatar(user.Id, &avatar)
	if err != nil {
		log.Printf("controller.AddAvatar ERROR failed to save avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, fmt.Sprintf("failed to save avatar[%s]", info.Filename), http.StatusBadRequest)
		return
	}
	avatar.Id = saved.Id
	tmpl := avatar.Template(user)
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("controller.AddAvatar ERROR failed to template avatar[%s], %s", info.Filename, err.Error())
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
