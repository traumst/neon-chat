package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"prplchat/src/db"
	"prplchat/src/handler"
	a "prplchat/src/model/app"
	"prplchat/src/model/event"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

const MaxUploadSize int64 = 10 * utils.KB

var allowedImageFormats = []string{
	"image/svg+xml",
	"image/jpeg",
	"image/gif",
	"image/png",
}

func isAllowedImageFormat(mime string) bool {
	for _, allowed := range allowedImageFormats {
		if allowed == mime {
			return true
		}
	}
	return false
}

func AddAvatar(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
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
		// &template.InfoMessage{
		// 	Header: "User is not authenticated",
		// 	Body:   "Your session has probably expired",
		// 	Footer: "Reload the page and try again",
		// }
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err = r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		log.Printf("[%s] AddAvatar ERROR multipart failed, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid multipart input"))
		return
	}
	file, info, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("[%s] AddAvatar ERROR reading input file failed, %s\n", reqId, err.Error())
		log.Printf("[%s] AddAvatar TRACE reading input file failed, %+v\n", reqId, info.Filename)
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
	mime := http.DetectContentType(fileBytes)
	if !isAllowedImageFormat(mime) {
		http.Error(w, "file type is not supported: "+mime, http.StatusBadRequest)
		return
	}
	oldAvatars, err := db.GetAvatars(user.Id)
	if err != nil {
		http.Error(w, "file type is not supported: "+mime, http.StatusBadRequest)
		return
	}
	saved, err := db.AddAvatar(user.Id, info.Filename, fileBytes, mime)
	if err != nil {
		log.Printf("controller.AddAvatar ERROR failed to save avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, fmt.Sprintf("failed to save avatar[%s]", info.Filename), http.StatusBadRequest)
		return
	}
	if len(oldAvatars) > 0 {
		for _, old := range oldAvatars {
			if old == nil {
				continue
			}
			err := db.DropAvatar(old.Id)
			if err != nil {
				log.Printf("controller.AddAvatar ERROR failed to drop old avatar[%v]", old)
			}
		}
	}
	avatar := handler.AvatarFromDB(*saved)
	tmpl := avatar.Template(user)
	html, err := tmpl.HTML()
	if err != nil {
		log.Printf("controller.AddAvatar ERROR failed to template avatar[%s], %s", info.Filename, err.Error())
		http.Error(w, fmt.Sprintf("failed to template avatar[%d]", avatar.Id), http.StatusBadRequest)
		return
	}
	if err = handler.DistributeAvatarChange(app, user, &avatar, event.AvatarChange); err != nil {
		log.Printf("controller.AddAvatar ERROR failed to distribute avatar[%s] update, %s", info.Filename, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func GetAvatar(app *handler.AppState, db *db.DBConn, w http.ResponseWriter, r *http.Request) {
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
		// &template.InfoMessage{
		// 	Header: "User is not authenticated",
		// 	Body:   "Your session has probably expired",
		// 	Footer: "Reload the page and try again",
		// }

		http.Redirect(w, r, "/", http.StatusFound)
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
