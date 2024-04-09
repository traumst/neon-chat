package controller

import (
	"fmt"
	"log"
	"net/http"

	"go.chat/utils"
)

func FavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "icon/favicon.ico")
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	//log.Printf("--%s-> ServeFile", utils.GetReqId(r))
	// _, err := utils.GetSessionCookie(r)
	// if err != nil {
	// 	utils.ClearSessionCookie(w)
	// 	log.Printf("--%s-> ServeFile WARN cookie\n", utils.GetReqId(r))
	// 	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	// 	return
	// }
	path := utils.ParseUrlPath(r)
	if len(path) < 1 {
		log.Printf("<-%s-- ServeFile ERROR args\n", utils.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filename := fmt.Sprintf("%s/%s", path[1], path[2])
	if filename == "" {
		log.Printf("<-%s-- ServeFile ERROR filename path[%v]\n", utils.GetReqId(r), filename)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//log.Printf("<-%s-- ServeFile TRACE serving %s\n", utils.GetReqId(r), filename)
	http.ServeFile(w, r, filename)
}
