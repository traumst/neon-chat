package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"neon-chat/src/consts"
)

func FavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/icon/scarab-bnw.svg")
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] requested [%s]\n", reqId, r.URL.Path)
	pathParts := strings.Split(r.URL.Path, "/")
	fileName := pathParts[len(pathParts)-1]
	ext := strings.Split(fileName, ".")

	var filePath string
	switch ext[len(ext)-1] {
	case "js":
		filePath = fmt.Sprintf("./static/script/%s", fileName)
	case "css":
		filePath = fmt.Sprintf("./static/css/%s", fileName)
	case "html":
		filePath = fmt.Sprintf("./static/html/%s", fileName)
	case "ico":
		filePath = fmt.Sprintf("./static/icon/%s", fileName)
	case "svg":
		filePath = fmt.Sprintf("./static/icon/%s", fileName)
	default:
		log.Printf("ERROR [%s] serving [%s]\n", reqId, r.URL.Path)
		w.Write([]byte("invalid path"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("TRACE [%s] served [%s]\n", reqId, filePath)
	http.ServeFile(w, r, filePath)
}
