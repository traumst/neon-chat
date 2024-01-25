package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.chat/utils"
)

func FavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "icons/favicon.ico")
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> ServeFile", GetReqId(r))
	_, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("--%s-> ServeFile WARN user, %s\n", GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	path := utils.ParseUrlPath(r)
	if len(path) < 1 {
		log.Printf("--%s-> ServeFile ERROR args\n", GetReqId(r))
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	filename := fmt.Sprintf("script/%s", path[2])
	if filename == "" {
		log.Printf("--%s-> ServeFile ERROR filename, %s\n", GetReqId(r), err)
		utils.SendBack(w, r, http.StatusBadRequest)
		return
	}

	log.Printf("--%s-> ServeFile TRACE serving %s\n", GetReqId(r), filename)
	http.ServeFile(w, r, filename)
}
