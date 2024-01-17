package controllers

import (
	"log"
	"net/http"

	"go.chat/utils"
)

func FavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "icons/favicon.ico")
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> ServeFile", reqId(r))
	_, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("--%s-> Home WARN user, %s\n", reqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	path := utils.ParseUrlPath(r)
	if len(path) < 1 {
		log.Printf("--%s-> ServeFile ERROR args\n", reqId(r))
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}

	filename := path[1]
	if filename == "" {
		log.Printf("--%s-> ServeFile ERROR filename, %s\n", reqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, filename)
}
